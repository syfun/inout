import React, { Component } from 'react'
import { Menu, Icon } from 'antd'

import { NavLink } from 'react-router-dom'

class LeftMenu extends Component {
  render () {
    return (
      <div style={{ width: 256 }}>
        <Menu
          defaultSelectedKeys={['1']}
          defaultOpenKeys={['sub1']}
          mode='inline'
          theme='dark'
          inlineCollapsed={false}
        >
          <Menu.Item key='1'>
            <NavLink to='/'>
              <Icon type='pie-chart' />
              <span>入库</span>
            </NavLink>

          </Menu.Item>
          <Menu.Item key='2'>
            <NavLink to='/pop'>
              <Icon type='desktop' />
              <span>出库</span>
            </NavLink>
          </Menu.Item>
          <Menu.Item key='3'>
            <NavLink to='/stock'>
              <Icon type='inbox' />
              <span>库存统计</span>
            </NavLink>
          </Menu.Item>
          <Menu.Item key='4'>
            <NavLink to='/item'>
              <Icon type='inbox' />
              <span>物品</span>
            </NavLink>
          </Menu.Item>
          <Menu.Item key='5'>
            <NavLink to='/setting'>
              <Icon type='inbox' />
              <span>设置</span>
            </NavLink>
          </Menu.Item>
        </Menu>
      </div>
    )
  }
}

export default LeftMenu
