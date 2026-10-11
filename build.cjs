// Empacotador opcional da plataforma. Requer Node.js apenas para manutenção, não para jogar.
const fs=require('fs'),path=require('path');const base=__dirname;
let template=fs.readFileSync(path.join(base,'source.html'),'utf8');
const css=fs.readFileSync(path.join(base,'styles.css'),'utf8');const js=fs.readFileSync(path.join(base,'engine3d.js'),'utf8')+'\n'+fs.readFileSync(path.join(base,'consent.js'),'utf8')+'\n'+fs.readFileSync(path.join(base,'app.js'),'utf8');
const favicon=fs.readFileSync(path.join(base,'assets/icon.svg')).toString('base64');
const build=standalone=>{
 let output=template.replace('<link rel="stylesheet" href="styles.css">','<style>\n'+css+'\n</style>');
 output=output.replace('<script src="engine3d.js"></script>\n<script src="consent.js"></script>\n<script src="app.js"></script>','<script>\n'+js+'\n</script>');
 output=output.replace('href="assets/icon.svg"','href="data:image/svg+xml;base64,'+favicon+'"');
 if(standalone)output=output.replace('<link rel="manifest" href="manifest.webmanifest">','');
 return output;
};
fs.writeFileSync(path.join(base,'index.html'),build(false));
fs.writeFileSync(path.join(base,'TwoSix_Play_STANDALONE.html'),build(true));
console.log('TwoSix Play 6.3: index.html e HTML independente gerados.');
