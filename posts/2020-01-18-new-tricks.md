I've been working on incrementally incorporating new tricks, techniques and
tools into my daily life to improve my workflow. These are some new things I've
recently begun using that I enjoy.

- [peco](https://github.com/peco/peco) Simplistic interactive filtering tool
- [yank](https://github.com/mptre/yank) Yank terminal output to clipboard
- `vim -` neovim/vim read from stdin
- `"*yy vim register`
- [hexyl](https://github.com/sharkdp/hexyl) A command-line hex viewer
- [bropages](http://bropages.org/) examples for command line programs

These CLI tools help further my OCD-driven quest to avoid using the mouse as
much as possible. Developers that work with the CLI clients for docker,
kubernetes, MySQL, etc often must take a value from the (often complex) output
of one command and use it as an input argument to another command. For example,
`docker ps` and `docker inspect` or `docker rm -f`. Or how about `kubectl
config get-contexts` and `kubectl config use-context "$context"`

### peco, yank, vim -, * vim register
By piping the output of a command to `peco` and `yank` you can filter results
and then use the HJKL keys to navigate values and copy a value to your
clipboard. This is much, much less frustrating than taking your hands off your
keyboard and going for a mouse, dragging the cursor over the value you want,
click+dragging to select it and then copying.

`docker`, `kubectl`, and other cli tools can often produce lengthy output in
your terminal. I've found that piping the outputs of these commands into VIM
where I'm comfortable searching and navigating is much better.

###### Two commands with long output. Piping the output into vim makes working with it easier
<pre class="prettyprint">
$ kubectl describe pod "$POD_ID" | vim -
$ docker inspect "$CID" | vim -
</pre>

###### See it in action
<script id="asciicast-I07WYu2kipYbEucbCrDT5jkKI" src="https://asciinema.org/a/I07WYu2kipYbEucbCrDT5jkKI.js" async></script>

### tmux
For a long time I used iterm2 panes and tabs to organize my terminals. My
primary editor is [neovim](https://neovim.io/), so I spend 95% of my time in my
terminal. I prefer to use a mouse as infrequently as possible so I bound ⌥
+hjkl to switch split panes. The biggest limitation I encountered on a daily
basis was I couldn't easily resize split panes without using my mouse. In the
interest of having a more cross-platform compatible development flow, I decided
to just switch to primarily tmux.
[.tmux.conf](https://github.com/cflynn07/dotfiles/blob/master/dots/.tmux.conf)

<img src="/static/images/Screen_Shot_2020-01-03_at_1.57.51_PM.png" alt="" style="width:100%"/>

### hexyl

### bropages
