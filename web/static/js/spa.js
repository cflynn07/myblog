/*
 * Just a super lazy script snippet to make my blog more SPA-like
*/

(() => {
  function fetchPath (path) {
    fetch(window.location.origin + path)
      .then((result) => result.text())
      .then((result) => {
        const newDocument = new DOMParser().parseFromString(result, 'text/html')
        const innerHTML = newDocument.querySelector('#maincontent').innerHTML
        document.querySelector('#maincontent').innerHTML = innerHTML
        PR.prettyPrint()
        window.scrollTo(0,0)
        bindEventLinks()
      })
  }

  function bindEventLinks () {
    document.querySelectorAll('nav.navbar a, #maincontent a').forEach((link) => {
      if (link.href.indexOf(window.location.origin) === -1) {
        // don't bind event to links to other domains
        return
      }
      link.addEventListener('click', (e) => {
        e.preventDefault()
        const path = e.currentTarget.attributes.href.value
        window.history.pushState({}, '', path)
        fetchPath(path)
      })
    })
  }

  document.addEventListener('DOMContentLoaded', () => {
    bindEventLinks()
    window.onpopstate = function (event) {
      fetchPath(window.location.pathname)
    }
  })
})()
