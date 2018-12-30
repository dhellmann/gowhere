# gowhere -- Tool for testing Apache redirect instructions

gowhere is a Golang re-implementation of
[whereto](https://docs.openstack.org/whereto/), and is used to test
the `redirect` and `redirectmatch` directives in an htaccess file to
detect incorrect redirection, cycles, and excessive hops.

## To-do list

- tests for processing functions
- pcre regexes?
