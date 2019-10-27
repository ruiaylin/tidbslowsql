import moment from 'moment'

const data = []
const current = moment()
const prev1Hour = current.subtract(1, 'hour').format('YYYY-MM-DD HH')
const prev2Hour = current.subtract(1, 'hour').format('YYYY-MM-DD HH')

for (let i = 0; i < 60; i++) {
  data.push({
    count: 0,
    date: new Date(`${prev2Hour}:${i < 10 ? `0${i}` : i}:00`)
  })
}

for (let i = 0; i < 60; i++) {
  data.push({
    count: 0,
    date: new Date(`${prev1Hour}:${i < 10 ? `0${i}` : i}:00`)
  })
}

export default data
