import React from 'react'
import AreaChart from '../AreaChart'
import { connect } from 'react-redux'

function renderPlan(plan, keys) {
  const header = keys
  const rows = plan.map(item => {
    const result = []
    header.forEach(h => {
      result.push(item[h])
    })
    return result
  })
  const maxWidths = rows.reduce((acc, row) => {
    return acc.map((item, i) => {
      return Math.max(item, row[i].length)
    })
  }, header.map(h => h.length))
  const renderSep = widths => '|' + widths.map(w => '-'.repeat(w + 2)).join('+') + '|'
  const pad = (s, length) => ' ' + s + ' '.repeat((length - s.length > 0 ? length - s.length : 0) + 1)
  const renderRow = (row, widths) => '|' + row.map((r, i) => pad(r, widths[i])).join('|') + '|'
  const render = (header, rows) =>
    [renderSep(maxWidths), renderRow(header, maxWidths)]
      .concat([renderSep(maxWidths)])
      .concat(rows.map(row => renderRow(row, maxWidths)))
      .concat([renderSep(maxWidths)])
      .join('\n') +
    '-'.repeat(maxWidths)

  return render(header, rows)
}

const Details = props => {
  let handleClose = () => {
    const app = document.querySelector('.tss-App')
    const details = document.querySelector('.tss-Details')

    app.style.transform = ''
    details.style.transform = 'translateX(-100vw)'
  }

  return (
    <article className="panel tss-Details">
      <div className="panel-heading tss-panel-heading">
        <div className="title is-4">SQL Info</div>
        <div onClick={handleClose}>
          <i
            style={{
              marginBottom: '1rem',
              cursor: 'pointer'
            }}
            className="fas fa-times-circle"
          />
        </div>
      </div>
      <div className="panel-block">
        <div className="title is-5">Trending</div>
      </div>
      <div className="panel-block">
        <AreaChart data={props.singleSqlTrending} />
      </div>
      <div className="panel-block sql-info">
        {props.sqlInfo &&
          props.sqlInfo.map(([k, v]) => {
            if (k === 'plan') {
              v = renderPlan(v, ['id', 'count', 'task', 'operator_info'])
              v = <pre>{v}</pre>
            }

            return (
              <div key={k}>
                <div className="title is-6">{k}</div>
                <div>{v}</div>
              </div>
            )
          })}
      </div>
    </article>
  )
}

const mapStateToProps = state => ({
  singleSqlTrending: state.singleSqlTrending.map(d => ({
    date: new Date(d.day_time),
    count: d.counter
  })),
  sqlInfo: state.sqlInfo
})

export default connect(
  mapStateToProps,
  null
)(Details)
