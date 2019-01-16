import router from '../../router'
import template from './index.html'

import './style.css'

export default class {
  mount(container) {
    document.title = 'foo'
    container.innerHTML = template
    container.querySelector('.foo__gobar').addEventListener('click', () => {
      router.go('/bar')
    })
  }
}