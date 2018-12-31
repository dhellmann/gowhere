# gowhere -- Tool for testing Apache redirect instructions

gowhere is a Golang re-implementation of
[whereto](https://docs.openstack.org/whereto/), and is used to test
the `redirect` and `redirectmatch` directives in an htaccess file to
detect incorrect redirection, cycles, and excessive hops.

## Example

The `example` directory contains 2 input files that can be used to
demonstrate how gowhere works:

    $ go run gowhere.go example/htaccess example/tests.txt
    Unexpected rule matched check on line 7: '/old_root/index.html' should produce 301 '/new_root/not_index.html'
        /old_root/index.html -> 301 /new_root/index.html [line 9]
    Cycle found from rule on line 11: '/cycle/a' should produce 301 '/cycle/d'
        /cycle/a -> 301 /cycle/b [line 11]
        /cycle/a -> 301 /cycle/c [line 12]
        /cycle/a -> 301 /cycle/a [line 13]
    Untested rule [line 4] redirect /project/def/new_page.html 301 /project/def/other_page.html
    Untested rule [line 7] redirectmatch ^/renamed/new1/ 301 /renamed/new2/
    Untested rule [line 12] redirect /cycle/b 301 /cycle/c
    Untested rule [line 13] redirect /cycle/c 301 /cycle/a
    
    2 failures
    exit status 1

## To-do list

- pcre regexes?
