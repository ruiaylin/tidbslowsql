import 'daterangepicker'

import React, { useEffect, useRef } from 'react'

import $ from 'jquery'
import moment from 'moment'

const DateTimePicker = props => {
  const inputEl = useRef(null)

  useEffect(() => {
    $(inputEl.current)
      .daterangepicker({
        timePicker: true,
        opens: 'left',
        applyButtonClasses: 'button is-primary is-small',
        cancelButtonClasses: 'button is-dark is-small',
        startDate: moment()
          .startOf('hour')
          .subtract(24, 'hour'),
        endDate: moment().startOf('hour'),
        locale: {
          format: 'M/DD hh:mm A'
        }
      })
      .on('apply.daterangepicker', (_, picker) => {
        const start = picker.startDate.format('YYYY-MM-DD HH:mm:ss')
        const end = picker.endDate.format('YYYY-MM-DD HH:mm:ss')

        if (typeof props.handlePickCallback === 'function') {
          props.handlePickCallback(start, end)
        }
      })

    if (typeof props.handlePickCallback === 'function') {
      props.handlePickCallback(
        $(inputEl.current)
          .data('daterangepicker')
          .startDate.format('YYYY-MM-DD HH:mm:ss'),
        $(inputEl.current)
          .data('daterangepicker')
          .endDate.format('YYYY-MM-DD HH:mm:ss')
      )
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <div className="field tss-DateTimePicker">
      <p className="control has-icons-left">
        <input
          style={{
            width: '300px',
            borderWidth: '2px'
          }}
          ref={inputEl}
          className="input"
          placeholder="Select Date Time Range"
        />
        <span className="icon is-small is-left">
          <i className="fas fa-calendar"></i>
        </span>
      </p>
    </div>
  )
}

export default DateTimePicker
