import React from 'react'
import { Tabs } from 'antd'

import Label from './Label'

const TabPane = Tabs.TabPane

const Setting = () => (
  <div>
    <Tabs defaultActiveKey='1'>
      <TabPane tab='物品分类' key='1'>
        <Label type='itemType' />
      </TabPane>
      <TabPane tab='仓库' key='2'>
        <Label type='warehouse' />
      </TabPane>
      <TabPane tab='人员' key='3'>
        <Label type='user' />
      </TabPane>
    </Tabs>
  </div>
)

export default Setting
