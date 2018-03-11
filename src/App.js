import React, { Component } from 'react'
import {
  BrowserRouter as Router,
  Route,
  Redirect
} from 'react-router-dom'

import './App.css'
import { Layout } from 'antd'
import LeftMenu from './components/LeftMenu'
import Setting from './components/Setting'
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
              <Redirect from='/' to='/push' />
              <Route path='/push' component={Push} />
              <Route path='/pop' component={Pop} />
              <Route path='/stock' component={Stock} />
              <Route path='/item' component={Item} />
              <Route path='/setting' component={Setting} />
            </div>
          </Content>
          <Footer style={{ textAlign: 'center' }}>
        仓库入库出库系统 ©2018 Created by syfun (sunyu418@gmail.com)
          </Footer>
        </Layout>
      </Router>
    )
  }
}

export default App
