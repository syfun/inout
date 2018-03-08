import React, { Component } from 'react'
import { Table, Input, Button, Popconfirm, InputNumber, Form, Modal } from 'antd'
import axios from 'axios'
import Selectx from './Selectx'

const FormItem = Form.Item

const ItemForm = Form.create()(
  (props) => {
    const { visible, onCancel, onOK, update, form } = props
    const { getFieldDecorator } = form
    const formItemLayout = {
      labelCol: { span: 4 },
      wrapperCol: { span: 14 }
    }
    return (
      <Modal
        visible={visible}
        title={update ? '更新物品信息' : '创建物品'}
        okText={update ? '更新' : '创建'}
        onCancel={onCancel}
        onOk={onOK}
        cancelText='取消'
      >
        <Form layout='horizontal'>
          <FormItem label='物品' {...formItemLayout}>
            {getFieldDecorator('name')(
              <Selectx
                url='/items'
                initialValue={update ? update.type : ''}
              />
            )}
          </FormItem>
          <FormItem label='物品分类' {...formItemLayout}>
            {getFieldDecorator('type')(
              <Selectx
                url='/labels?type=itemType'
                initialValue={update ? update.type : ''}
              />
            )}
          </FormItem>
          <FormItem label='物品规格' {...formItemLayout}>
            {getFieldDecorator('specification', {initialValue: update ? update.specification : ''})(
              <Input />
            )}
          </FormItem>
          <FormItem label='单位' {...formItemLayout}>
            {getFieldDecorator('unit', {initialValue: update ? update.unit : ''})(
              <Input />
            )}
          </FormItem>
          <FormItem label='入库' {...formItemLayout}>
            {getFieldDecorator('push', {initialValue: update ? update.push : 0})(
              <InputNumber min={0} />
            )}
          </FormItem>
          <FormItem label='出库' {...formItemLayout}>
            {getFieldDecorator('pop', {initialValue: update ? update.pop : 0})(
              <InputNumber min={0} />
            )}
          </FormItem>
          <FormItem label='剩余' {...formItemLayout}>
            {getFieldDecorator('now', {initialValue: update ? update.now : 0})(
              <InputNumber min={0} />
            )}
          </FormItem>
        </Form>
      </Modal>
    )
  }
)

class Item extends Component {
  state = {
    data: [],
    visible: false,
    update: null
  }

  columns = [
    {
      title: '物品名称',
      dataIndex: 'name'
    }, {
      title: '物品分类',
      dataIndex: 'type'
    }, {
      title: '物品规格',
      dataIndex: 'specification'
    }, {
      title: '单位',
      dataIndex: 'unit'
    }, {
      title: '摘要',
      dataIndex: 'abstract'
    }, {
      title: '入库数量',
      dateIndex: 'number'
    }, {
      title: '仓库',
      dateIndex: 'warehouse'
    }, {
      title: '入库人',
      dateIndex: 'user'
    }, {
      title: '备注',
      dateIndex: 'remark'
    }, {
      title: '操作',
      render: (text, record) => {
        return (
          <div>
            <Button
              type='primary'
              style={{ marginRight: '5px' }}
              onClick={() => this.setState({visible: true, update: record})}
            >
              编辑
            </Button>
            <Popconfirm
              placement='bottomLeft'
              title='确定删除吗？'
              onConfirm={() => this.handleDelete(record.id)}
              okText='确定' cancelText='取消'
            >
              <Button>删除</Button>
            </Popconfirm>
          </div>
        )
      }
    }
  ]

  componentDidMount () {
    axios.get('/items').then(res => {
      this.setState({data: res.data})
      this.cacheData = res.data.map(item => ({ ...item }))
    })
  }

  handleCancel = () => {
    this.setState({visible: false})
  }

  handleCreate = () => {
    const form = this.form
    form.validateFields((_, values) => {
      axios.post('/items', values).then(
        res => {
          let data = [...this.state.data]
          data.push(res.data)
          this.setState({ data, visible: false })
        }
      )
    })
  }

  handleUpdate = () => {
    const form = this.form
    const itemID = this.state.update.id
    form.validateFields((_, values) => {
      axios.patch(`/items/${itemID}`, values).then(
        res => {
          const data = [...this.state.data]
          const target = data.filter(item => itemID === item.id)[0]
          Object.assign(target, res.data)
          this.setState({ data, visible: false })
        }
      )
    })
  }

  handleDelete = (key) => {
    axios.delete(`/items/${key}`).then(
      res => {
        let data = [...this.state.data]
        data = data.filter(item => key !== item.id)
        this.setState({ data, visible: false })
      }
    )
  }

  showModal = () => {
    this.setState({visible: true})
  }

  saveFormRef = (form) => {
    this.form = form
  }

  render () {
    const {update, visible} = this.state
    return (
      <div>
        <Button
          className='editable-add-btn'
          onClick={() => this.setState({visible: true, update: null})}
        >
          添加
        </Button>
        <ItemForm
          ref={this.saveFormRef}
          visible={visible}
          onCancel={this.handleCancel}
          onOK={update ? this.handleUpdate : this.handleCreate}
          update={update}
        />
        <Table
          bordered
          dataSource={this.state.data}
          columns={this.columns}
          rowKey='id'
        />
      </div>
    )
  }
}

export default Item
