# gowhere -- Tool for testing Apache redirect instructions

gowhere is a Golang re-implementation of
[whereto](https://docs.openstack.org/whereto/), and is used to test
the `redirect` and `redirectmatch` directives in an htaccess file to
detect incorrect redirection, cycles, and excessive hops.

## Example

The `example` directory contains 2 input files that can be used to
demonstrate how gowhere works:

    $ go get github.com/dhellmann/gowhere/cmd/gowhere

    $ cd $GOPATH/src/github.com/dhellmann/gowhere

    $ $GOPATH/bin/gowhere example/htaccess example/tests.txt
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

    $ echo $?
    1

## Inputs

To test a set of redirects, `gowhere` needs the input `.htaccess`
file and another input file with test data.

The `.htaccess` file should contain `Redirect` and
`RedirectMatch` directives. Blank lines and lines starting with
`#` are ignored. For example, this input includes 6 rules:

    # Redirect old top-level HTML pages to the version under most recent
    # full release.
    redirectmatch 301 ^/$ /current-release/
    redirectmatch 301 ^/index.html$ /current-release/
    redirectmatch 301 ^/projects.html$ /current-release/projects.html
    redirectmatch 301 ^/language-bindings.html$ /current-release/language-bindings.html

    # Redirect subpage pointers to main page
    redirect 301 /install/ /current-release/install/
    redirect 301 /basic-install/ /current-release/install/

    # this is gone and never coming back, indicate that to the end users
    redirect 410 /obsolete_content.html

The test data file should include one test per line, including 3
parts: the input path, the expected HTTP response code, and the
(optional) expected output path. For example:

    / 301 /current-release/
    / 301 /current-release
    /install/ 301 /current-release/install/
    /no/rule 301 /should/fail
    /obsolete-content.html 410

    # verify that this path is not redirected
    /current-release/index.html 200


## To-do list

- pcre regexes?
