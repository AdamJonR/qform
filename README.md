# QForm

QForm is a DSL for creating HTML5 forms. It is implemented using Dialects, a recursive-descent parser for Domain Specific Languages (DSLs) that is implemented using Go and facilitates parsing through use of Parsing Expression Grammars (PEGs).

## Motivation

I get tired of writing and sifting through HTML, but I especially lament the redundant, superfluous extravagances required by the markup of HTML5 forms.

## Examples

A few examples of the basic usage of QForm are listed below with their corresponding outputs.

### A Basic Contact Form

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

```
<form method="post">
  <div class="form-group">
    <label for="name">Name</label>
    <input type="text" name="name" maxlength="30" required="required" id="name" />
  </div>
  <div class="form-group">
    <label for="email">Email</label>
    <input type="email" name="email" id="email" />
  </div>
  <div class="form-group">
    <label for="message">Message</label>
    <textarea name="message" id="message"></textarea>
  </div>
  <div class="form-group">
    <input type="submit" value="Send message" name="field4" id="field4" />
  </div>
</form>
```

## Using QForm Directly

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

## Using QForm With Polyglot

You can automatically parse files with the QForm DSL using [Polyglot](https://github.com/AdamJonR/polyglot), a program that packages multiple DSL parsers together and allows you parse several different DSLs within the same file.
