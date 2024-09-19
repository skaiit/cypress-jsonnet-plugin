const os = require('os');
const fs = require('fs');
const { execSync } = require('child_process');
const path = require('path');
 
const jsonnetModulePath = path.join('./node_modules/cy-jsonnet/jsonnet');
 
let commandToRun = "";
 
/**
* Build code to generate the jsonnet binary with custom extensions
* Set the commandToRun so other function can run it for jsonnet interpretation
*/
function buildPlugin() {
  try {
    console.info('Generating cy-jsonnet binary');
    execSync(`go build -C ${jsonnetModulePath}`);
    let runner = path.join(jsonnetModulePath, 'cy-jsonnet');
    commandToRun = `${runner}`;
    // If windows OS, setup accordingly
    if (os.platform().includes('win32')) {
      runner = path.join(jsonnetModulePath, 'cy-jsonnet.exe');
      commandToRun = `cmd /c ${runner}`;
    }
    else {
      console.log('Setting the shell script permission to be executable');
      execSync(`chmod +x ${runner}`);
    }
    console.info('Command to run = '+commandToRun);
    console.info('Generating cy-jsonnet binary');
  } catch (err) {
    console.error('An error occurred during jsonnet extension build');
    console.error(err);
    process.exit(1);
  }
}
 
/**
* interpretJsonnet provides support to generate the json
* @param {string} jsonnetRootFolder - jsonnet root folder for all template files
* @param {string} fileSearchPattern - jsonnet file name pattern e.g. '*.jsonnet'
* @param {string} outputFolder -path for output json
* @param {string} generateTestData -Std.extVar("generateTestData") to avoid any jsonnet->json
*/
function interpretJsonnet(jsonnetRootFolder,fileSearchPattern, outputFolder, generateTestData = false) {
  try {
    let command = `${commandToRun} --jsonnetRootFolder=${jsonnetRootFolder} --fileSearchPattern=${fileSearchPattern} --outputFolder=${outputFolder}`
    if (generateTestData){
      command = command + " --generateTestData"
    }
    const commandShellStdout = execSync(command).toString();
    console.log(commandShellStdout);
  } catch (err) {
    console.error('An error occurred during jsonnet evaluation for dynamic test creation:', jsonnetRootFolder);
    console.error(err);
    process.exit(1);
  }
}
 
buildPlugin();
module.exports = {
  interpretJsonnet,
}; 