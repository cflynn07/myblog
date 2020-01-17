<img src="/static/images/linux-command-line-shell-scripting-bible.png" style="width:50%" />

I'm focusing on reading more technical books this year. Recently I finished
"Linux Command Line and Shell Scripting Bible (3rd Edition)." I've been working
with linux and shell scripting for 15+ years and pretty much everything I know
I learned through some amalgamation of blog posts and informal unstructured
reading. I chose to read this book to try in fill in any gaps in my knowledge
that I may have.

The book is about 700 pages and it took me about 10 days to finish it. Some
sections were too simple to really be informative, but in the interest of being
thorough I read them. Thankfully there were plenty of sections containing
information that was new to me and I definitely picked up some new tricks.

The book focuses on linux utilities found in the GNU coreutils package, which
differ slightly from their counterparts found in BSD descended systems (such as
OS X). I primarily work in OS X, so when I started reading this book I
downloaded an ubuntu image and VirtualBox to have a proper linux VM to work in.
I also used brew on OS X to replace most of my BSD utilities with their GNU
counterparts. (`brew install coreutils`)

Some sections of the book are just too simple to be of much value to anyone
with moderate experience, such as Ch.3 which covers basic commands like mv, rm,
ls, cd, etc. In the interest of being thorough, I made sure not to skip these
sections just in case there was some tidbit of info I didn't know.

By going into detail about various components of linux, the book increased the
depth of my knowledge on subjects I was at least moderately knowledgeable in.
For example, it's been a long time since I've managed physical and logical
volumes, partitions, filesystems, user accounts and groups. Using a VM with
VirtualBox made trying the logical volume examples easy, since I could create
virtual disks and attach/mount them on my VM.

New things I picked up include:
  - using `shift` with input and function parameters
  - opening, closing and some redirecting of file descriptors
  - trapping process signals
  - standard cli tool options parsing with `getopt` and `getopts`

There were actually a few things that I actually knew and was surprised to see
weren't covered. Such as:
##### bash parameter expansion.
<pre class="prettyprint">
# Default value for a variable
$ echo ${foo-defaultval}
defaultval

# How to get a substring?
$ foo="welcome to the fold"

# First 5 chars of $foo
$ echo ${foo:0:5}
welco

# how many characters in $foo
echo ${#foo}
20

# Last 5 chars of $foo
$ echo ${foo:${#foo}-5:${#foo}}
fold.
</pre>

##### The "::" library naming convention
<pre class="prettyprint">
ExternalModule::utilityFunction
</pre>

##### !!
<pre class="prettyprint">
$ whoami
casey

$ !!
casey
</pre>

##### xargs
This I've used for years for things like:
<pre class="prettyprint">
$ docker ps -q | xargs docker rm -f
$ cat ./utils
ps
man
git
$ cat names | xargs where
/bin/ps
/usr/bin/man
git: aliased to hub
/usr/bin/git
</pre>

##### tr
<pre class="prettyprint">
# easier to read PATH dirs listing
$ echo $PATH | tr ':' '\n' | sort
/Applications/MySQLWorkbench.app/Contents/MacOS
/Library/Frameworks/Python.framework/Versions/3.8/bin
/Users/casey/.antigen/bundles/Dbz/kube-aliases
/Users/casey/.antigen/bundles/bhilburn/powerlevel9k
/Users/casey/.antigen/bundles/robbyrussell/oh-my-zsh/lib
/Users/casey/.antigen/bundles/robbyrussell/oh-my-zsh/plugins/command-not-found
/Users/casey/.antigen/bundles/robbyrussell/oh-my-zsh/plugins/git
/Users/casey/.antigen/bundles/robbyrussell/oh-my-zsh/plugins/lein
/Users/casey/.antigen/bundles/robbyrussell/oh-my-zsh/plugins/pip
/Users/casey/.antigen/bundles/zsh-users/zsh-autosuggestions
/Users/casey/.antigen/bundles/zsh-users/zsh-syntax-highlighting
/Users/casey/.cargo/bin
/Users/casey/go/bin
/Users/casey/google-cloud-sdk/bin
/bin
/sbin
/usr/bin
/usr/local/bin
/usr/local/go/bin
/usr/sbin
</pre>

The last third of the book focuses on `awk` and `sed`. This definitely a deeper
dive into either of those utilities than I've ever done before. My feelings are
that awk is the more useful of the two and it can be used effectively with
extracting substrings from command output for piping to additional commands.
Ex:
<pre class="prettyprint">
function cinspect {
  docker inspect $(docker ps | awk "/$1/{print \$1}")
}
$ cinspect "blog" | head -n 5
[
    {
        "Id": "2b0380143b40b7f3cc7d14665553552cc12da1567e244407cf15142a9f2c7b32",
        "Created": "2020-01-17T07:52:54.0597444Z",
        "Path": "/bin/sh",

what's the docker name of a container?
$ docker ps | awk '/blog/{print $NF}'
k8s_blog_blog-app-5fb57645cf-54h7k_blog_5d82ae3b-38fe-11ea-b74b-025000000001_0
</pre>

I struggle to think of situations I typically encounter where I'd use some of
the advanced awk/sed abilities that are covered, such as the sed `n` or `N`
next commands. The seed's been planted in my head that these things exist,
maybe I'll come up with something.

For anyone who's considering reading through this book I'd recommend using
[shellcheck](https://github.com/koalaman/shellcheck) as a linter to check your
scripts. It's definitely useful for catching syntax errors and pointing out
best practices.

After practicing writing shell scripts for a few consecutive days, I'm more
likely to identify potential repetitive tasks that I can turn into shell
scripts to improve my productivity (I make a point to try and add these to my
[dotfiles](https://github.com/cflynn07/dotfiles) repo). I also used shell
scripts in projects where previously I probably would have written a js/nodejs
or golang script.
([cedict-sqlite3](https://github.com/cflynn07/cedict-sqlite3))
