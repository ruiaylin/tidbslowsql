import { make as Container } from './Container.bs'
import { make as Navbar } from './components/Navbar/Navbar.bs'
import React from 'react'
import { Provider } from 'react-redux'
import configureStore from './store'
import Details from './components/Details'

const store = configureStore()

const App = () => (
  <Provider store={store}>
    <div className="tss-App">
      <Navbar />
      <Container />
    </div>
    <Details />
  </Provider>
)

export default App
