weblowercaser
=============

A small program to change all file names and internal links to lowercase.

Why
---

Because you inherited a huge website made on Windows with incoherent case.

Compile it
----------

	go build

If it fails, you may need to run first

	go get golang.org/x/net/html

Use it
------

... at your own risks. It seems to work. And it's fast. But seriously I had tested it only for my sites, it might need some complements for your own. It doesn't handle PHP files and generally can't change incomplete HTML files. Links you build using JavaScript won't be fixed either.

This makes a fixed copy in destpath of the source site :

	./weblowercaser -from sourcepath -to destpath


