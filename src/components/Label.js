import React, { Component } from 'react'
import { Table, Input, Button, Popconfirm, Form, Modal } from 'antd'
import axios from 'axios'

const FormItem = Form.Item

const ItemTypeForm = Form.create()(
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
        title={update ? '更新' : '创建'}
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
        </Form>
      </Modal>
    )
  }
)

class ItemType extends Component {
  state = {
    data: [],
    visible: false,
    update: null
  }

  columns = [
    {
      title: '名称',
      dataIndex: 'name'
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
    axios.get(`/labels?type=${this.props.type}`).then(res => {
      this.setState({data: res.data || []})
    })
  }

  handleCancel = () => {
    this.setState({visible: false})
  }

  handleCreate = () => {
    const form = this.form
    form.validateFields((_, values) => {
      values.type = this.props.type
      axios.post('/labels', values).then(
        res => {
          let data = [...this.state.data]
          data.push(res.data)
          form.resetFields()
          this.setState({ data, visible: false })
        }
      )
    })
  }

  handleUpdate = () => {
    const form = this.form
    const itemID = this.state.update.id
    form.validateFields((_, values) => {
      axios.patch(`/labels/${itemID}`, values).then(
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
    axios.delete(`/labels/${key}`).then(
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
        <ItemTypeForm
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

export default ItemType
