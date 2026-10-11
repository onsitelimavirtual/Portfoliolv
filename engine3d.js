/* TwoSix Play — WebGL 3D nativo, sem biblioteca remota. Geometria low-poly. */
(()=>{'use strict';
function init(canvas){
 const gl=canvas?.getContext('webgl',{alpha:false,antialias:true,preserveDrawingBuffer:false});
 if(!gl)return null;
 const vertex=`attribute vec3 p;attribute vec3 n;attribute vec3 c;varying vec3 v;uniform float fog;void main(){vec3 q=p-vec3(0.,3.,-7.);float z=q.z*.994-q.y*.11;float y=q.y*.994+q.z*.11;gl_Position=vec4(q.x*2.5,y*2.5,z-.2,z);float light=max(.28,dot(normalize(n),normalize(vec3(-.35,.85,-.45))));float haze=clamp((z-7.)/150.,0.,.7);v=mix(c*light,vec3(.32,.47,.68),haze);}`;
 const fragment=`precision mediump float;varying vec3 v;void main(){gl_FragColor=vec4(v,1.);}`;
 const shader=(type,src)=>{const s=gl.createShader(type);gl.shaderSource(s,src);gl.compileShader(s);if(!gl.getShaderParameter(s,gl.COMPILE_STATUS))throw Error('Shader 3D inválido');return s;};
 let prog;try{prog=gl.createProgram();gl.attachShader(prog,shader(gl.VERTEX_SHADER,vertex));gl.attachShader(prog,shader(gl.FRAGMENT_SHADER,fragment));gl.linkProgram(prog);if(!gl.getProgramParameter(prog,gl.LINK_STATUS))throw Error('Programa 3D indisponível');}catch(_){return null;}
 gl.useProgram(prog);const buf=gl.createBuffer();gl.bindBuffer(gl.ARRAY_BUFFER,buf);for(const [name,offset] of [['p',0],['n',3],['c',6]]){const a=gl.getAttribLocation(prog,name);gl.enableVertexAttribArray(a);gl.vertexAttribPointer(a,3,gl.FLOAT,false,36,offset*4);}gl.enable(gl.DEPTH_TEST);gl.disable(gl.CULL_FACE);
 let out=[];
 function tri(a,b,c,color){const u=[b[0]-a[0],b[1]-a[1],b[2]-a[2]],v=[c[0]-a[0],c[1]-a[1],c[2]-a[2]];let n=[u[1]*v[2]-u[2]*v[1],u[2]*v[0]-u[0]*v[2],u[0]*v[1]-u[1]*v[0]];const l=Math.hypot(...n)||1;n=n.map(x=>x/l);for(const pt of [a,b,c])out.push(...pt,...n,...color);}
 function quad(a,b,c,d,color){tri(a,b,c,color);tri(a,c,d,color);}
 function box(x,y,z,w,h,d,color){const x0=x-w/2,x1=x+w/2,y0=y,y1=y+h,z0=z-d/2,z1=z+d/2;
 quad([x0,y0,z0],[x1,y0,z0],[x1,y1,z0],[x0,y1,z0],color);
 quad([x0,y0,z1],[x0,y1,z1],[x1,y1,z1],[x1,y0,z1],color);
 quad([x0,y0,z0],[x0,y1,z0],[x0,y1,z1],[x0,y0,z1],color);
 quad([x1,y0,z0],[x1,y0,z1],[x1,y1,z1],[x1,y1,z0],color);
 quad([x0,y1,z0],[x1,y1,z0],[x1,y1,z1],[x0,y1,z1],color);
 }
 function cone(x,y,z,r,h,color){for(let i=0;i<7;i++){const a=i*Math.PI*2/7,b=(i+1)*Math.PI*2/7;tri([x+Math.cos(a)*r,y,z+Math.sin(a)*r],[x,y+h,z],[x+Math.cos(b)*r,y,z+Math.sin(b)*r],color);}}
 function wheel(x,y,z,r,color){for(let i=0;i<10;i++){const a=i*Math.PI*2/10,b=(i+1)*Math.PI*2/10;const ca=Math.cos(a),sa=Math.sin(a),cb=Math.cos(b),sb=Math.sin(b);quad([x-.16,y+ca*r,z+sa*r],[x+.16,y+ca*r,z+sa*r],[x+.16,y+cb*r,z+sb*r],[x-.16,y+cb*r,z+sb*r],color);}}
 function vehicle(x,z,kind,style=0,color){ if(kind==='bonus'){cone(x,1.0,z,.65,1.5,[1,.75,.22]);return;}
 if(kind==='moto'){wheel(x-.78,.44,z-.92,.42,[.08,.10,.16]);wheel(x+.82,.44,z+.92,.42,[.08,.10,.16]);box(x,.82,z,.58,.22,1.55,color||[.97,.37,.18]);box(x-.05,1.05,z-.12,.48,.24,.72,[.12,.21,.33]);box(x+.48,1.22,z-.18,.26,.52,.28,[.16,.18,.24]);box(x-.62,1.45,z-.05,.28,.28,.28,[.99,.74,.22]);box(x-.08,.98,z-.62,.18,.10,.55,[.83,.88,.92]);return;}
 const colors=color||[.10,.85,.93]; box(x,.23,z,1.72,.44,3.28,colors); box(x,.66,z-.12,1.46,.58,1.56,[.06,.16,.29]); if(style===1||style===4) box(x,.42,z+1.02,1.54,.30,1.02,[.13,.18,.26]); if(style===2||style===5) box(x,.98,z-.06,1.12,.12,1.05,[.24,.83,.98]); else box(x,.96,z-.42,1.10,.10,1.02,[.22,.78,.96]); for(const dx of [-1,1])for(const dz of [-1,1])box(x+dx*.84,.12,z+dz*1.01,.36,.56,.68,[.06,.08,.16]); }
 function render(d,type,time){out=[];const race=type==='race';gl.viewport(0,0,gl.drawingBufferWidth,gl.drawingBufferHeight);gl.clearColor(race?.045:.20,race?.08:.35,race?.19:.38,1);gl.clear(gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT);
 const c=race?[.14,.19,.27]:[.30,.30,.18];box(0,-.24,77,120,.20,180,race?[.1,.29,.39]:[.28,.38,.27]);box(0,-.06,75,10,.13,180,c);
 if(race){for(let i=0;i<11;i++){let z=(((i*19-d.dist*2)%190)+190)%190+3;for(const lane of [-1.5,1.5])box(lane,.03,z,.11,.02,7,[.7,.92,1]);for(const side of [-5.4,5.4])box(side,.03,z,.10,.06,8,[.13,.89,.97]);}}
 else {for(let i=0;i<20;i++){let z=(((i*10-d.dist*1.3)%160)+160)%160+2;for(const side of [-7,7]){box(side,0,z,.38,.6,.37,[.43,.34,.20]);cone(side,.6,z,2.1,3.4,[.14,.52,.35]);}}}
 for(let i=0;i<18;i++){const z=10+i*14;const side=i%2?1:-1;cone(side*(race?17:24),-.2,z,7,(i%4)*1.8+5,race?[.24,.35,.51]:[.24,.43,.40]);}
 for(const o of d.objects.filter(o=>!o.claimed&&o.dist>=d.dist&&o.dist-d.dist<130)){
 const z=8+(o.dist-d.dist)*.8;
 if(race){if(o.kind==='bonus')cone(o.lane*3.2,0,z,.8,2,[1,.74,.17]);else vehicle(o.lane*3.2,z,'car',o.shape,[[.37,.78,.99],[.95,.42,.55],[.55,.62,.99],[.65,.88,.60],[.95,.72,.34],[.24,.89,.84]][o.color%6]);}
 else{if(o.kind==='bonus')cone(o.lane*3.2,1.0,z,.75,1.2,[1,.70,.14]);else if(o.kind==='log') box(o.lane*3.2,.10,z,1.4,.52,2.1,[.54,.36,.24]); else box(o.lane*3.2,.10,z,1.2,.92,1.1,[.47,.38,.29]);}
 }
 if(race){vehicle(d.lane*3.2,4,'car',3,[.20,.85,.96]);}
 else {vehicle(d.lane*3.2,4,'moto',0,[.98,.35,.18]);}
 gl.bufferData(gl.ARRAY_BUFFER,new Float32Array(out),gl.DYNAMIC_DRAW);gl.drawArrays(gl.TRIANGLES,0,out.length/9);
 }
 return {render};
}
window.TwoSix3D={init};
})();
