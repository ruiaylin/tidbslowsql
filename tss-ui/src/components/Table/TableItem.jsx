import React from 'react'
import moment from 'moment'

const TableItem = props => {
  const item = props.item

  const handleShowDetails = () => {
    props.handleGetSingleSqlTrending(item.cluster, item.sql_id)
    props.handleGetSqlInfo(item.cluster, item.query_id)

    const app = document.querySelector('.tss-App')
    const details = document.querySelector('.tss-Details')

    app.style.transform = 'translateX(50vw)'
    details.style.transform = 'translateX(10vw)'
  }

  return (
    <tr
      style={{
        background: item.affirm === 1 ? '#eaeaea' : ''
      }}
      className="tss-TableItem"
      onClick={handleShowDetails}
    >
      <td>
        {item.cluster} {item.affirm === 1 && <i className="fas fa-check has-text-primary"></i>}
      </td>
      {item.inst_name && <td>{item.inst_name}</td>}
      <td>{item.schema_name}</td>
      {item.count && <td>{item.count}</td>}
      <td>{item.sql_id}</td>
      <td className="sql-text" data-tip={item.sql_text} data-for="sqlText">
        {item.sql_text && item.sql_text.substring(0, 50) + '...'}
      </td>
      {item.day_time && (
        <td className="datetime">
          {moment(
            moment(item.day_time)
              .utc()
              .format('YYYY-MM-DD HH:mm:ss')
          ).fromNow()}
        </td>
      )}
      {item.first_seen && (
        <td className="datetime">
          {moment(
            moment(item.first_seen)
              .utc()
              .format('YYYY-MM-DD HH:mm:ss')
          ).fromNow()}
        </td>
      )}
      <td>{item.query_time.toFixed(3)}</td>
      <td>{item.process_keys}</td>
      <td>{item.process_time.toFixed(3)}</td>
      <td>{item.commit_time.toFixed(3)}</td>
      <td>{item.total_keys}</td>
      <td>{item.write_keys}</td>
      <td className="max-memory">{(item.mem_max / 1024).toFixed(0) + ' KB'}</td>
    </tr>
  )
}

export default TableItem
