import React, { Component } from 'react'
import { Table, Input, Button, Popconfirm, InputNumber, Form, Modal, Row, Col } from 'antd'
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
          <FormItem label='名称' {...formItemLayout}>
            {getFieldDecorator('name', {initialValue: update ? update.name : ''})(
              <Input />
            )}
          </FormItem>
          <FormItem label='物品分类' {...formItemLayout}>
            {getFieldDecorator('type', {initialValue: update ? update.type : ''})(
              <Selectx
                url='/labels?type=itemType'
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
    update: {},
    pagination: {},
    filter: {}
  }

  columns = [
    {
      title: '物品名称',
      dataIndex: 'name',
      width: '15%'
    }, {
      title: '物品分类',
      dataIndex: 'type',
      width: '15%'
    }, {
      title: '物品规格',
      dataIndex: 'specification',
      width: '15%'
    }, {
      title: '单位',
      dataIndex: 'unit',
      width: '10%'
    }, {
      title: '入库',
      dataIndex: 'push',
      width: '10%'
    }, {
      title: '出库',
      dataIndex: 'pop',
      width: '10%'
    }, {
      title: '剩余',
      dataIndex: 'now',
      width: '10%'
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
  }

  fetch = () => {
    axios.get('/items?page=1&page_size=10', {
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
    axios.get(`/items?page=${pager.current}&page_size=${pager.pageSize}`).then(
      res => {
        this.setState({
          data: res.data || [],
          pagination: pager
        })
      }
    )
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
          form.resetFields()
          let pager = this.state.pagination
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
      axios.patch(`/items/${itemID}`, values).then(
        res => {
          const data = [...this.state.data]
          const target = data.filter(item => itemID === item.id)[0]
          Object.assign(target, res.data)
          this.setState({ data, visible: false })
          form.resetFields()
        }
      )
    })
  }

  handleDelete = (key) => {
    axios.delete(`/items/${key}`).then(
      res => {
        let data = [...this.state.data]
        data = data.filter(item => key !== item.id)
        const pager = this.state.pagination
        pager.total -= 1
        this.setState({ data, visible: false, pagination: pager })
      }
    )
  }

  showModal = () => {
    this.setState({visible: true})
  }

  saveFormRef = (form) => {
    this.form = form
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

  render () {
    const {update, visible, pagination} = this.state
    return (
      <div>
        <Row>
          <Col span={4} offset={16}>
            <Input
              placeholder='物品名称'
              onChange={(e) => this.setFilter('name', e.target.value)}
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
          pagination={pagination}
          onChange={this.handlePageChange}
        />
      </div>
    )
  }
}

export default Item
