const yMiliSecond = 31536e+6
const mMiliSecond = 2629746e+3
const dMilisecond = 8.64e+7

function calculate() {
    const currentDate = new Date()
    const birthday = getBirthday()
    let result = currentDate - birthday

    let rYear = Math.trunc(result/yMiliSecond)
    result -= rYear*yMiliSecond
    let rMonth = Math.trunc(result/mMiliSecond)
    result -= rMonth*mMiliSecond
    let rDay = Math.trunc(result/dMilisecond)

    let timeOfLife = {year: rYear, month: rMonth, days: rDay}
    buildResult(timeOfLife)
}

function getBirthday() {
    const day = () => isNaN(parseInt(document.getElementById('day').value)) ? 0 : parseInt(document.getElementById('day').value);
    const month = () => isNaN(parseInt(document.getElementById('month').value)) ? 0 : parseInt(document.getElementById('month').value);
    const year = () => isNaN(parseInt(document.getElementById('year').value)) ? 0 : parseInt(document.getElementById('year').value);

    return new Date(year(), month()-1, day())
}

function buildResult({year: y, month: m, days: d}) {
    document.getElementById('r-year').innerHTML = `${y} years`
    document.getElementById('r-month').innerHTML = `${m} months`
    document.getElementById('r-days').innerHTML = `${d} days`
}