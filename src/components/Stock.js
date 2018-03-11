import React, { Component } from 'react'
import { Table, Input, Button, Popconfirm, InputNumber, Form, Modal, Select, Row, Col } from 'antd'
import axios from 'axios'

import Selectx from './Selectx'

const FormItem = Form.Item
const Option = Select.Option

const StockForm = Form.create()(
  (props) => {
    const { visible, onCancel, onOK, update, form, handleSelectChange, items } = props
    const { getFieldDecorator } = form
    const formItemLayout = {
      labelCol: { span: 4 },
      wrapperCol: { span: 14 }
    }
    return (
      <Modal
        visible={visible}
        title={update ? '更新' : '创建'}
        okText={update ? '更新' : '创建'}
        onCancel={onCancel}
        onOk={onOK}
        cancelText='取消'
      >
        <Form layout='horizontal'>
          <FormItem label='物品' {...formItemLayout}>
            {getFieldDecorator('name', {
              initialValue: update ? update.name : '',
              rules: [{
                required: true, message: '请选择物品！'
              }]
            })(
              <Select onChange={handleSelectChange}>
                {
                  items.map(item => <Option key={item.id}>{item.name}</Option>)
                }
              </Select>
            )}
          </FormItem>
          <FormItem label='物品分类' {...formItemLayout}>
            {getFieldDecorator('type', {initialValue: update ? update.type : ''})(
              <Input />
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
          <FormItem label='仓库' {...formItemLayout}>
            {getFieldDecorator('warehouse', {
              initialValue: update ? update.warehouse : '',
              rules: [{
                required: true, message: '请选择仓库！'
              }]
            })(
              <Selectx
                url='/labels?type=warehouse'
              />
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
          <FormItem label='item_id' style={{display: 'none'}}>
            {getFieldDecorator('item_id')(
              <Input />
            )}
          </FormItem>
        </Form>
      </Modal>
    )
  }
)

class Stock extends Component {
  state = {
    data: [],
    visible: false,
    update: null,
    items: [],
    pagination: {},
    filter: {}
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
      title: '入库',
      dataIndex: 'push'
    }, {
      title: '出库',
      dataIndex: 'pop'
    }, {
      title: '剩余',
      dataIndex: 'now'
    }, {
      title: '仓库',
      dataIndex: 'warehouse'
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
    this.fetch()

    axios.get('/items').then(res => {
      this.setState({items: res.data || []})
    })
  }

  fetch = () => {
    axios.get('/stocks?page=1&page_size=10', {
      params: this.state.filter
    }).then(res => {
      this.setState({
        data: res.data || [],
        pagination: {
          pageSize: 10,
          total: parseInt(res.headers.total, 10),
          showTotal: total => `共 ${res.headers.total} 条`
        }
      })
    })
  }

  handlePageChange = (pagination, filters, sorter) => {
    const pager = this.state.pagination
    pager.current = pagination.current
    axios.get(`/stocks?page=${pager.current}&page_size=${pager.pageSize}`).then(
      res => {
        this.setState({
          data: res.data || [],
          pagination: pager
        })
      }
    )
  }

  handleSelectChange = (value) => {
    const form = this.form
    const {items} = this.state
    const target = items.filter(item => item.id === parseInt(value, 10))[0]
    form.setFieldsValue({
      type: target.type,
      unit: target.unit,
      specification: target.specification,
      item_id: target.id
    })
  }

  handleCancel = () => {
    this.setState({visible: false})
  }

  handleCreate = () => {
    const form = this.form
    form.validateFields((_, values) => {
      console.log(values)
      axios.post('/stocks', values).then(
        res => {
          let data = [...this.state.data]
          data.push(res.data)
          form.resetFields()
          const pager = this.state.pagination
          pager.total += 1
          this.setState({ data, visible: false, pagination: pager })
        }
      )
    })
  }

  handleUpdate = () => {
    const form = this.form
    const itemID = this.state.update.id
    form.validateFields((_, values) => {
      axios.patch(`/stocks/${itemID}`, values).then(
        res => {
          const data = [...this.state.data]
          const target = data.filter(item => itemID === item.id)[0]
          Object.assign(target, res.data)
          form.resetFields()
          this.setState({ data, visible: false })
        }
      )
    })
  }

  handleDelete = (key) => {
    axios.delete(`/stocks/${key}`).then(
      res => {
        let data = [...this.state.data]
        data = data.filter(item => key !== item.id)
        this.setState({ data, visible: false })
      }
    )
  }

  setFilter = (key, value) => {
    const filter = this.state.filter
    if (value === '') {
      delete filter[key]
    } else {
      Object.assign(filter, {[key]: value})
    }
    this.setState(filter)
  }

  showModal = () => {
    this.setState({visible: true})
  }

  saveFormRef = (form) => {
    this.form = form
  }

  render () {
    const {update, visible, items, pagination} = this.state
    return (
      <div>
        <Row>
          <Col span={3} offset={8}>
            <Input
              placeholder='物品名称'
              onChange={(e) => this.setFilter('name', e.target.value)}
              onPressEnter={() => this.fetch()}
            />
          </Col>
          <Col span={3} offset={1}>
            <Input
              placeholder='物品分类'
              onChange={(e) => this.setFilter('type', e.target.value)}
              onPressEnter={() => this.fetch()}
            />
          </Col>
          <Col span={3} offset={1}>
            <Input
              placeholder='仓库'
              onChange={(e) => this.setFilter('warehouse', e.target.value)}
              onPressEnter={() => this.fetch()}
            />
          </Col>
          <Col span={2}>
            <Button
              className='editable-add-btn'
              onClick={this.fetch}
            >
          搜索
            </Button>
          </Col>
          <Col span={2}>
            <Button
              className='editable-add-btn'
              onClick={() => this.setState({visible: true, update: null})}
            >
          添加
            </Button>
          </Col>
        </Row>
        <StockForm
          ref={this.saveFormRef}
          visible={visible}
          onCancel={this.handleCancel}
          onOK={update ? this.handleUpdate : this.handleCreate}
          update={update}
          handleSelectChange={this.handleSelectChange}
          items={items}
        />
        <Table
          bordered
          dataSource={this.state.data}
          columns={this.columns}
          rowKey='id'
          pagination={pagination}
          onChange={this.handlePageChange}
        />
      </div>
    )
  }
}

export default Stock
