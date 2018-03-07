import React, { Component } from 'react'
import axios from 'axios'
import { Select } from 'antd'
const Option = Select.Option

class Selectx extends Component {
  constructor (props) {
    super(props)
    this.state = {
      data: []
    }
  }
  componentDidMount () {
    axios.get(this.props.url).then(
      res => this.setState({data: res.data})
    )
  }
  render () {
    const {data} = this.state
    return (
      <div>
        <Select defaultValue={this.props.initialValue} onChange={this.props.onChange}>
          {
            data.map(item => (
              <Option key={item.name}>{item.name}</Option>
            ))
          }
        </Select>
      </div>
    )
  }
}

export default Selectx
