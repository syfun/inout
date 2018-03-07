import React, { Component } from 'react'
import {
  BrowserRouter as Router,
  Route
} from 'react-router-dom'

import './App.css'
import { Layout } from 'antd'
import LeftMenu from './components/LeftMenu'
import Label from './components/Label'
import Push from './components/Push'
import Pop from './components/Pop'
import Stock from './components/Stock'
import Item from './components/Item'

const { Content, Footer, Header } = Layout

class App extends Component {
  render () {
    return (
      <Router>
        <Layout className='layout'>
          <Header>
            <LeftMenu />
          </Header>
          <Content style={{ padding: '0 50px', marginTop: '20px' }}>
            <div style={{ padding: 24, background: '#fff', minHeight: 280, textAlign: 'right' }}>
              <Route exact path='/' component={Push} />
              <Route path='/pop' component={Pop} />
              <Route path='/stock' component={Stock} />
              <Route path='/item' component={Item} />
              <Route path='/setting' component={Label} />
            </div>
          </Content>
          <Footer style={{ textAlign: 'center' }}>
        Ant Design ©2016 Created by Ant UED
          </Footer>
        </Layout>
      </Router>
    )
  }
}

export default App
