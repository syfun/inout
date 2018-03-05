import React, { Component } from 'react'
import axios from 'axios'
import { Select } from 'antd'
const Option = Select.Option

class Selectx extends Component {
  constructor (props) {
    super(props)
    const { url, defaultValue } = props
    this.state = {
      data: [],
      url,
      defaultValue
    }
  }
  componentDidMount () {
    axios.get(this.state.url).then(
      res => this.setState({data: res.data})
    )
  }
  render () {
    const {data, defaultValue} = this.state
    return (
      <div>
        <Select defaultValue={defaultValue} style={{ width: 120 }}>
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
