//spinner
let spinner = `<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span><span class="visually-hidden"> Давай подождём...</span>`;

// DateTimeToString
function DateTimeToString(datetime) {
    let year = datetime.getUTCFullYear();
    let month = datetime.getUTCMonth()+1;
    let day = datetime.getUTCDate();
    let hour = datetime.getUTCHours();
    let minute = datetime.getUTCMinutes();
    let second = datetime.getUTCSeconds();
    return `${day<10?'0'+day:day}.${month<10?'0'+month:month}.${year} ${hour<10?'0'+hour:hour}:${minute<10?'0'+minute:minute}:${second<10?'0'+second:second}`
}

// Date Now
function NowDateTime(){
    let now = new Date();
    let year = now.getFullYear();
    let month = now.getMonth()+1;
    let day = now.getDate();
    $('#sdate').val(`${year}-${month<10?'0'+month:month}-${day<10?'0'+day:day}T00:00`);
    $('#edate').val(`${year}-${month<10?'0'+month:month}-${day<10?'0'+day:day}T23:59`);
}

// Timer
let timerCount = 0;
let timerEnabled = false;
function startTimer(id) {
    timerEnabled = true;
    $(`#${id}`).text(timerCount + " сек");
    setTimeout(function () {
        timerCount++;
        $(`#${id}`).text(timerCount + " сек");
        if (timerEnabled) {
            startTimer(id);
        } else {
            timerCount = 0;
        }
    }, 1000);
}

function stopTimer() {
    timerEnabled = false;
}

// Random colors
function getRandomColor() {
    let r = Math.round(Math.random() * (255 - 150) + 100);
    let g = Math.round(Math.random() * (255 - 150) + 100);
    let b = Math.round(Math.random() * (255 - 150) + 100);
    return `rgb(${r}, ${g}, ${b});`;
}

//DataTables
function setDatatables(id, rows, is_excel) {
    let buttons = [];
    if (is_excel) {
        buttons.push({
            extend: 'excelHtml5',
            text: 'Excel',
        })
    }
    $(`#${id}`).DataTable({
        data: rows,
        destroy: true,
        dom: 'Bfrtip',
        buttons: buttons,
    });
}