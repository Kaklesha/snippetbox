# **ui\html\pages** section info 

### Displaying dynamic data

### Dynamic content escaping
The html/template package automatically escapes any data that is yielded between {{ }}
tags. This behavior is hugely helpful in avoiding cross-site scripting (XSS) attacks, and is the
reason that you should use the html/template package instead of the more generic
text/template package that Go also provides.
As an example of escaping, if the dynamic data you wanted to yield was

 ```
<span>{{"<script>alert('xss attack')</script>"}}</span>
 ```
It would be rendered harmlessly as:

```
<span>&lt;script&gt;alert(&#39;xss attack&#39;)&lt;/script&gt;</span>
```

The html/template package is also smart enough to make escaping context-dependent. It
will use the appropriate escape sequences depending on whether the data is rendered in a
part of the page that contains HTML, CSS, Javascript or a URI.

  ``` ```