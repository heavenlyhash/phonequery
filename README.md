phonequery
==========

A jank scraper to collect phone stats and filter 'em.

Run `./run` to run.

Patch it to get different answers.  Ain't no fancy config here.  Grep for `[]filter{` for where all the filters are declared.

First time it's run, it saves a bunch of json to a cache file.  This is slow; many HTTP, such not parallel.  Delete the cache if you want to be slow again.


