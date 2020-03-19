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

# p81
- bottom level components known as leaf components

# p82
- https://reactjs.org/docs/thinking-in-react.html
  - also references single responsibility principle
  - build top down or bottom up
  - if passed in from parent via props not state?

# p83
- forms specia state managers (stateful situation even when properties are passed down from parent)
- sometimes make parent components just for holding state

# p84
- TimersDashboard logical home of state not EditableTimerList because of need to create

# p91
- avoid initializing state of input field to `undefined`
- introduction of input fields initialized by state becoming out of sync

# p92
- React onChange attribute
- combination of `value` and `onChange` attributes is how form elements are handled in react

# p96
- noticing naming strategy
  - name of method passed down to child: `onSomethingSomething` (comes in on props)
  - name of custom component method that "calls up" `handleSomethingSomethin`
  - EX: <foo onSomethingSomething={this.handleSomethingSomething}></foo>

# p108
- `forceUpdate` introduced
- `componentDidMount` `componentWillUnmount`

# p112
- another example of 'wiring up' event handling propagation from child to parent components

# p122
```
NOW="$(date +%s)000"
# UUID=$(uuidgen | tr '[:upper:]' '[:lower:]')
UUID="a73c1d19-f32d-4aff-b470-cea4e792406a"
DATA=$(printf '{"start":%s,"id":"%s"}' $NOW $UUID)
curl -X POST \
  -H 'Content-Type: application/json' \
  -d "$DATA" \
  http://localhost:3000/api/timers/start
```

# p126
- since I was last paying attention, AJAX handled by "fetch" object/api in browsers

# p130
- JS Date object doesn't stringify as presented in book?

# p135
- The virtual DOM introduced, tree of React elements
- virtual dom vs shadow dom
- https://www.webcomponents.org/community/articles/introduction-to-shadow-dom
- https://developers.google.com/web/fundamentals/web-components/shadowdom?hl=en

# p141
- introduction of ReactText object

# p143
- Explanation of react parser and jsx (javascript syntax extension)

# p146
- JSX attributes need values, ex:
`<input name='foo' disabled={true} />` {true} required

# p147
- spread syntax 好用 <Component {...props} />

# p148
- recommendation npm package `classnames`
- for/htmlFor gotcha
- html entitites/emoji "gotcha"

# p150
- html element attributes must be prefixed with `data-`

# p151
- aria accessability attributes
- TODO: read links on p151

# p153
- contrast of ReactComponent vs ReactElement, ReactComponent.render() returns ReactElement

# p157
- `context` introduced "implicit props"

# p158
- `propTypes` introduced

# p159
- proptypes can validate simple scalar types or do more complex validations
- `defaultProps`

# p160
- pass context to all kids: `childContextTypes` and `getChildContext`

# p162
- example of `childContextTypes` and `getChildContext()` in a component

# p163
- `contextTypes` in child element tells react what context properties the child wants
so in short:
PARENT: [childContextTypes{} + getChildContext()] --> CHILD: [contextTypes{}]

# p164
- first example of require() on css file

# p165
- functional stateless components
- context global potential good use case: logged in user

# p169
- render function onClick value calls a func that returns a func

# p171
- webpack CSS encapsulation

# p172
- stateful components required class property `state`

# p173
- remember constructor run once and before component mounted

# p174
- *pass function to setState(), when state update depends on current state, preferable to pass function
- setState asynchronous
- example of user spamming the decrement button faster that state asynch updates
```
this.setState((prevState) => {
  return {
    value: prevState.value - 1
  }
})
```

# p176
- mitigate/minimize complex states build apps single stateful component composed of stateless components

# p177
- extracting some functionality of Switch to stateless component Choice

# p179
- props.children, sort of 'transposing' elements
EX:
<Container>STUFF</Container
class Container extends React.Component {
  render () {
    return (
      <div className='container'>{this.props.children}</div>
    )
  }
}
^ "STUFF" will be wrapped in div.container

# p180
- propTyes.oneOf()
static propTypes = {
  children: PropTypes.oneOf([...])
}

# p181
- React.Children.map()
- React.Children.forEach()
- React.cloneElement()
- React.createElement()

# p182
- example of "wrapping" each element in a list of child elements with a wrapper component
  - makes use of React.Children.map()

# p183
- React.Children.toArray()

# p184
- coming in next ch, lifecycle methods like `componentWillUpdate()` - useful for stuff like form validation

# p188
- SyntheticMouseEvent && .nativeEvent (cross browser standardization wrapper)

# p189
- onMouseMove, (other event handler props listed) - also other groups of events (ex keyboard, focus, animation, transition, etc)

# p191
- shared event handler func for button (using event object to determine which button clicked)

# p192
- using 'refs' property to access DOM elements in a component (useful in forms for getting values)

# p195
- must use `key` property when children in iterator (read Dynamic Children docs)

# p196
- controlled v/ uncontrolled components (text field accessed via event object) (react is not controlling the input field, therefore uncontrolled)

# p197
- converting form elements to controlled has advantages, validation, localstorage for persisting 1/2 completed
