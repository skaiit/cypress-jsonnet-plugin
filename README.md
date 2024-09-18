
<h1 align="center">Welcome to Cypress-Jsonnet plugin 👋</h1>

## ✨ Description

`cy-jsonnet` plugin is designed to auto generate cypress test definition and payload using jsonnet templates.

**Core Ideas :** 
1. Eliminate duplication with object-orientation
1. Keep common data models in jsonnet file
1. Share and extend the models according to your user needs
1. Balance the logic between jsonnet and Cypress specs, e.g keep simple branching of positive and negative test cases only
1. **Few things to note:**
Jsonnet is hermetic: It always generates the same data no matter the execution environment. The addition of the native function was to support random data so we can generate idempotent payload for each test run.

Dynamic test case generation can be advantageous and must have native support for randomizing specific properties for the test payload. Jsonnet offers a solution to this by providing support for native methods extension, allowing for customization according to individual requirements.

## Quick overview about Jsonnet
[Jsonnet](https://jsonnet.org/) is a configuration language for application and tool developers developed by Google.

| Core ideas    | ... |
| -------- | ------- |
| Generate config data  | Open source (Apache 2.0)    |
| Side-effect free | Familiar syntax     |
| Organize, simplify, unify    | Reformatter, linter    |
| Manage sprawling config | Editor & IDE integrations     |
| simple extension of JSON    | Formally specified    |


## 🚀 How to use this plugin for Cypress testing?

1. Pre-requisite : Make sure you have these softwares 
   - [Golang](https://go.dev/doc/install)
   - [node js](https://nodejs.org/en/download)
   - [VS Code](https://code.visualstudio.com/download)
   - [Go-Jsonnet](https://github.com/google/go-jsonnet) 
      Run `go install github.com/google/go-jsonnet/cmd/jsonnet@latest`


2. Create folder where you like to test this plugin.
3. Open your favorite terminal window and cd to the folder
4. Init Node project `npm init -y`
5. Install Cypress `npm install cypress`
6. Install TypeScript `npm install -g typescript`
7. Install this plugin `npm install cy-jsonnet`
8. Open Cypress `npx cypress open`
   - This will create required cypress folders and config files
9. Add tsconfig.json to project `tsc --init`
10. Open VS Code `code .`
    - Edit `cypress.config.ts`

```TS
import { interpretJsonnet } from "cy-jsonnet";
import { defineConfig } from "cypress";
import * as path from 'path';
const jsonnetFolder: string = "jsonnet";
const testDefinitionFolder: string = "testDefinition";
const testDataFolder: string = "testData";
let jsonnetPath = path.join('./cypress/support/', jsonnetFolder);
let testDefinitionPath = path.join('./cypress/fixtures/', testDefinitionFolder);
let testDataPath = path.join('./cypress/fixtures/', testDataFolder);
export default defineConfig({
  e2e: {
    setupNodeEvents(on, config) {
      on("before:run", (details) => {
        interpretJsonnet(jsonnetPath, '**/*.jsonnet', testDefinitionPath, false);
      });
      on("before:spec", (spec) => {
        interpretJsonnet(jsonnetPath, `**/*${spec.fileName}*.jsonnet`, testDataPath, true);
      });
      return config;
    },
    experimentalRunAllSpecs: true,
    experimentalInteractiveRunEvents: true
  },
});
```
   - Edit `tsconfig.json`

```TS
{
  "compilerOptions": {
    "target": "es2021",
    "lib": ["es2021", "dom"],
    "moduleResolution": "node",
    "resolveJsonModule": true,
    "types": ["cypress", "node", "@cypress/grep"],
    "strictNullChecks": true,
    "allowSyntheticDefaultImports": true,
    "baseUrl": "./",
    "paths": {
      "@fixtures/*": ["cypress/fixtures/*"],
      "@support/*": ["cypress/support/*"]
    }
  },
  "include": ["**/*.ts", "**/*.js"], 
}
```

10. Now your Cypress test setup is ready to configure dynamic test example using this plugin.
   - Create a new file "cypress\e2e\person.cy.ts"
```TS
import * as testInfo from "@fixtures/testDefinition/person.json";
before(() => {
  cy.fixture("testData/person.json").as("testData");
});
describe("Cypress Dynamic Tests Positive Scenario Examples", () => {
  testInfo.positiveScenarios.forEach((testMetaData, index) => {
    it(testMetaData.testDefinition.scenario, { tags: testMetaData.testDefinition.tags }, function (data: any = this.testData.positiveScenarios[index].testData) {
      expect(true).equals(true);
    // Add your test logic here
      console.log(data.person);
    });
  });
});
describe("Cypress Dynamic Tests Negative Scenario Examples", () => {
  testInfo.negativeScenarios.forEach((testMetaData, index) => {
    it(testMetaData.testDefinition.scenario, { tags: testMetaData.testDefinition.tags }, function (data: any = this.testData.negativeScenarios[index].testData) {
      expect(true).equals(true);
 // Add your test logic here
      console.log(data.person);
    });
  });
}); 
```

11. Add required jsonnet and libsonnet files for using standard pattern
   - Add libsonnet file "cypress\support\jsonnet\lib\utils.libsonnet"
```JS
{
  testDefinition(fileName, scenario, tags=[])::
    {
      scenario: scenario,
      testIdentifier: std.md5(fileName+scenario),
      tags: tags + [self.testIdentifier],
    },
}

```
   - Add libsonnet file "cypress\support\jsonnet\lib\models.libsonnet"

```JS
{
  Person(firstName=std.native("fake")("{firstname}"),lastName=std.native("fake")("{lastname}"), ssn=std.native("fake")("{ssn}"), address={})::
    {
      person: {
        firstName: firstName,
        lastName: lastName,
        ssn:ssn,
        address: address,
      },
    },
  Address(city='bellevue', state=std.native('fake')('{state}'))::
    {
      city: city,
      state: state,
    },
}
```
   - Add sample jsonnet files "cypress\support\jsonnet\person.jsonnet"
```JS
local model = import './lib/models.libsonnet';
local utils = import './lib/utils.libsonnet';
local DynamicTest(definition={}, data={}) = {
  testDefinition: definition,
  testData: if std.extVar('generateTestData') == 'true' then data else {},
};
{
  positiveScenarios: [
    // Person
    DynamicTest(
      definition=utils.testDefinition(fileName=std.thisFile,scenario='Verify person entity can be created with empty address', tags=['sanity', 'person']),
      data=model.Person()
    ),
    DynamicTest(
      definition=utils.testDefinition(fileName=std.thisFile,scenario='Verify person is created with name=joe', tags=['sanity', 'regression']),
      data=model.Person(firstName='joe')
    )
  ],
  negativeScenarios: [
    DynamicTest(
      definition=utils.testDefinition(fileName=std.thisFile,scenario='Verify person failed with empty firstname', tags=['regression', 'sanity']),
      data=model.Person(firstName="",address=model.Address(city=std.native("fake")("{city}")))
    ),
    DynamicTest(
      definition=utils.testDefinition(fileName=std.thisFile,scenario='Verify person failed when ssn in null or empty', tags=['regression']),
      data=model.Person(ssn="")
    ),
  ],
}
```

12. Now we can see how dynamic test and data from jsonnet can load into person.cy.ts
   - Run 🏃 `npx cypress open --e2e`
   - Select `person.cy.ts`
13. To use Cypress grep plugin
   - Install `npm install @cypress/grep`
   - Add these lines to `cypress\support\e2e.ts`

```TS
import registerCypressGrep from '@cypress/grep'
registerCypressGrep()
```

14. This will help with running using grep tags
         ` npx cypress run --env grepTags="sanity"`
    - More examples [@cypress/grep](https://www.npmjs.com/package/@cypress/grep)

### 📜 Use all the supported features of [jsonnet](https://jsonnet.org/learning/tutorial.html) and [gofakeit](https://pkg.go.dev/github.com/brianvoe/gofakeit/v7#readme-simple-usage).


## 🤝 Contribution

Contributions, issues and feature requests are welcome. Please email your ideas to us.<br />

## Authors

👤 **Sumit Agarwal**   [BD :email:](sumit.agarwal@bd.com) [Personal :email:](ska4siva@gmail.com) 

👤 **Anoop Sasi**  [BD :email:](anoop.sasi@bd.com) [Personal :email:](anoopsasi@outlook.com)

## Code Contributors

This project exists thanks to all the people who contributed.
- Kunal Nain [BD :email:](kunal.nain@bd.com)
- David Lukic-Hanomihl [📧](david@assertqa.com)
- Marko Kolasinac [📧](marko.kolasinac@assertqa.com)
- 
## Show your support

If you find this useful please spread the words :thumbsup:

## 📝 License

This project is [MIT] licensed.

## Special thanks 🙏 
- go-jsonnet
- golang
- node js
- Cypress
- gofakeit
- cypress/grep

---

