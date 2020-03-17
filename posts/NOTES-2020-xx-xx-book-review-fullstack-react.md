# p2
- JS representation of DOM, diffing in userland

# p16
- http//semantic-ui.com (bootstrap alternative)

# p19
- first annoying side quest, figure out how to make `standardjs` compatible with project
  - added standard config to package.json

# p20
- comparison of JSX vs JS (React.createElement())
- JSX is compiled into javascript

# p23
- babel introduced, transpling
  - for now using "on-the-fly" -> will explore transpiling for production later

# p26
- first mention of react-dom

# p27
- first mention of child/parent component relationships

# p28
- 'class' reserved word, hence we use className. 

# p31
- first mention of `props`
- braces are delmiters

# p35
- first example using Array.prototype.map()
- first example using `key` attribute/property

# p39
- "a child does not own its props" - one way data flow
- passing down functions in props is "the canonical manner"

# p41
- onClick handler first mentioned
- `this` context binding. Custom component methods contexts null
  - React binds context automatically only for defaut set of API methods

# p42
- review ES5/ES6 classes https://gist.github.com/remarkablemark/fa62af0a2c57f5ef54226cae2258b38d

# p43
- "state" introduced

# p44
- this.setState() introduced

# p45
- contrast of components not owning props but owning state

# p46
- react lifecycle methods introduced

# p47
- this.setState() required for state modification

# p49
- first mention of javascript pass by reference
- Array.prototype.concat() > push(), concat doesn't mutate

# p53
- babel plugin: transform-class-properties
- introduction of babel plugins and presets
- babel-standalone
  - default uses 2 presets:  (preset set of plugins used to support particular lang features)
    - es2015
    - react

# p54
- javascript features stages (1-4)
- https://github.com/tc39/proposal-class-public-fields
  - https://github.com/tc39/proposal-class-fields (stage 3)
- had to add to package.json:
```
"standard": {
  "parser": "babel-eslint",
  "globals": [
    "React",
    "ReactDOM"
  ]
}
```
and install babel-eslint (--save-dev)
side read: https://itnext.io/property-initializers-what-why-and-how-to-use-it-5615210474a3
- method definition order matters

# p61
- new project steps to work with standardjs
  - npm install babel-eslint eslint --save-dev
  - add ^ standard prop to package.json
- skipping jsx-a11y/label-has-for (TODO read)

# p66
- https://en.wikipedia.org/wiki/Single-responsibility_principle
- https://blog.cleancoder.com/uncle-bob/2014/05/08/SingleReponsibilityPrinciple.html
  - "We want to increase the cohesion between things that change for the same
    reasons, and we want to decrease the coupling between those things that
    change for different reasons.
- TimersList & TimersDashboard components

# p67
- ToggleableTimerForm component

# p69
- EditableTimer component, child either Timer or TimerForm
- remame TimerList to EditableTimerList

# p71
- layout of steps to build app from scratch
  - good first step, define heirarchy, then build static version
- in this ch example top component will talk to server (TimersDashboard)

# p77
- react property defaultValue
