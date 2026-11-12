package main

import "html/template"

var page = template.Must(template.New("report").Parse(`<!doctype html>
<html lang="en"><head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>RFC 20 · ASCII interchange lab</title>
<style>
:root { background: #edf3f3; color: #152b35; font: 17px/1.5 system-ui, sans-serif; }
main { max-width: 900px; margin: auto; padding: 36px 22px; }
h1 { font-size: clamp(2.1rem, 6vw, 3.6rem); line-height: 1.08; margin: .3em 0; }
p { max-width: 70ch; }
a { color: #115e67; }
code { overflow-wrap: anywhere; }
.receivers { display: grid; grid-template-columns: repeat(auto-fit, minmax(min(100%, 280px), 1fr)); gap: 14px; }
section { background: white; border: 1px solid #aac1c0; border-radius: 10px; padding: 16px 20px; }
h2 { margin: 0; font-size: 1.3rem; }
dt { font-weight: 700; margin-top: .6em; }
dd { margin-left: 0; }
ol { margin: 0; padding-left: 1.5em; }
</style></head><body><main>
<h1>RFC 20 · Same ASCII, different record boundaries</h1>
<p>All {{len .Input}} input bytes have their high bit clear. RFC 20 describes two receivers that use different break characters.</p>
<p><strong>Input</strong> <code>{{printf "%q" .Input}}</code><br><strong>Hex bytes</strong> <code>{{printf "% X" .Input}}</code></p>
<div class="receivers">
{{range .Receivers}}
<section aria-label="{{.Name}} receiver">
<h2>{{.Name}}</h2>
<dl><dt>Break character</dt><dd><code>{{printf "%q" .Break}}</code></dd>
<dt>Complete records</dt><dd>{{if .Records}}<ol>{{range .Records}}<li><code>{{printf "%q" .}}</code></li>{{end}}</ol>{{else}}None{{end}}</dd>
<dt>Unterminated tail</dt><dd>{{if .Tail}}<code>{{printf "%q" .Tail}}</code>{{else}}None{{end}}</dd></dl>
</section>{{end}}
</div>
<p>This saved report comes from the Go CLI. Quoted strings show control characters as escapes.</p>
<p><a href="https://www.rfc-editor.org/info/rfc20/">Official RFC 20</a> · <a href="https://github.com/Sreenivas-Sadhu-Prabhakara/rfc-0020-ascii-lab">Source and run instructions</a></p>
</main></body></html>
`))
