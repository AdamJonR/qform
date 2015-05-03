QForm
=====

QForm is a DSL for creating HTML5 forms. It is implemented using Dialects, a recursive-descent parser for Domain Specific Languages (DSLs) that is implemented using Go and facilitates parsing through use of Parsing Expression Grammars (PEGs).

Motivation
----------

I get tired of writing and sifting through HTML, but I especially lament the redundant, superfluous extravagances required by the markup of HTML5 forms.

Examples
--------

Below is a basic contact form:

```
- method post

text
- name name
- maxlength 30
- required

email
- name email

textarea
- name message

submit
- value Send message
```

Below are some options fields:

```
radio
- name preference
- [
  call Call me back
  email Email me a message
  mail Send me a letter
]

checkbox
- name permission
- [
  yes I give my permission to contact me
]

select
- name department
- [
  sales
  tech Tech Support
  receivables
]
```



Using QForm Directly
--------------------

You can implement a custom application in Go and use the QForm library directly. A generalized usage example is provided in the code below.
```
...
// store file contents
bytes, err := ioutil.ReadFile(inputPath)
// check for error
if err != nil {
  // exit early if read error
  return
}
// convert to string
source := string(bytes)
// create dsl
dsl := new(qform.DSL)
// parse source
output, err = dialects.Parse(dsl, source)
// check for parsing error
if err != nil {
  // handle as is best for your app
}
// convert output to bytes[]
byteSource := []byte(output)
// try to save source to output path
err = ioutil.WriteFile(outputPath, byteSource, 0777)
...
```

Using QForm With Polyglot
-------------------------

You can automatically parse files with the QForm DSL using [Polyglot](https://github.com/AdamJonR/polyglot), a program that packages multiple DSL parsers together and allows you parse several different DSLs within the same file.
