import React, { useState, useEffect } from 'react'
import ReactTooltip from 'react-tooltip'

const Table = props => {
  const Header = props.header
  const Item = props.item
  const [asc, setAsc] = useState(false)
  const [tableData, setTableData] = useState([])

  useEffect(() => {
    setTableData(props.data)
  }, [props.data])

  const handleSortTable = (f1, f2) => () => {
    if (asc) {
      setTableData(tableData.sort(f2))
      setAsc(false)
    } else {
      setTableData(tableData.sort(f1))
      setAsc(true)
    }
  }

  return (
    <>
      <div className="tss-TableWrapper">
        <table className="table is-fullwidth is-bordered is-hoverable tss-Table">
          <Header handleSortTable={handleSortTable} item={props.data.length > 0 ? props.data[0] : {}} />
          <tbody>
            {props.data.length > 0 &&
              props.data.map((item, i) => (
                <Item
                  key={item.query_id + item.sql_id + i}
                  item={item}
                  handleGetSqlInfo={props.handleGetSqlInfo}
                  handleGetSingleSqlTrending={props.handleGetSingleSqlTrending}
                />
              ))}
          </tbody>
        </table>
      </div>
      {props.data.length === 0 && (
        <div className="subtitle has-text-centered table-no-data">No data available in table</div>
      )}
      {tableData.length > 0 && (
        <ReactTooltip className="tss-Tooltip" id="sqlText" getContent={datatip => <div>{datatip}</div>} />
      )}
    </>
  )
}

export default Table
