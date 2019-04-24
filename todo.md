# To do ...

An unordered thinking space for new features.

## 3 Next big things

1. Keep (capped) list of requests and responses
2. 
3. 


## Increase adoption

 - powershell auto complete
 - go releaser
    - release versioning 
    - chocolety
    - homebrew
    - apt-get
 - import from
   - swagger
   - cli history
   - postman
   - curl statement

## Random

- be able to interact with the history
- what-if => dump url cmd or json
- read in json (or last cmd) as input values
- convert value to type within payload e.g.
  - `"age":"{int:arg}"` => `"age":5`
  - `"isEnabled":"{bool:arg}"` => `"isEnabled":false` using standard truthy/falsy
 - defaultQuery for an api, many apis take the api token as a query parameter
 - separate query args in cmd (from path)
 - does it already work for XML or json to XML conversion?
- download binary
- payload a binary file
- form posting

 ## Work around already exists
 
 - dynamic functions for parameters e.g. (all currently doable outside of capi, using the normal shell but might help)
   - `{fn:uuid}` - v4 uuid
   - `{fn:date(yyyyMMdd)}` - current date (UTC only) formatted
     - yyMMdd, etc.
     - enum (go consts, Kitchen, etc.)
     - unix second/ms/nano
   - `{fn:dateRel(fmt,-1h)}` relative time, an hour ago
   - `{fn:rndstr(5)}` - random string, 5 chars long
 - full blown templating engine `{{ if .value }}foo{{ end }}`
- override profile,
  - headers
  - query string
  - payload
