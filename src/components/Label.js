import React from 'react'
import { Tabs } from 'antd'

import Tablex from './Tablex'

const TabPane = Tabs.TabPane

const columns = [
  {
    title: '名称',
    dataIndex: 'name',
    width: '30%'
  }
]

const Label = () => (
  <div>
    <Tabs defaultActiveKey='1'>
      <TabPane tab='物品分类' key='1'>
        <Tablex
          columns={columns}
          resource='label'
          query='type=itemType'
          extraData={{type: 'itemType'}}
        />
      </TabPane>
      <TabPane tab='仓库' key='2'>
        <Tablex
          columns={columns}
          resource='label'
          query='type=warehouse'
          extraData={{type: 'warehouse'}}
        />
      </TabPane>
      <TabPane tab='人员' key='3'>
        <Tablex
          columns={columns}
          resource='label'
          query='type=user'
          extraData={{type: 'user'}}
        />
      </TabPane>
    </Tabs>
  </div>
)

export default Label
