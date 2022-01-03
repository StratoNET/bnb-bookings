let attention = Inform();

(function () {
  'use strict';
  window.addEventListener('load', function () {
    // fetch all the forms needed to apply custom Bootstrap validation styles
    let forms = document.getElementsByClassName('needs-validation');
    // loop through them & prevent submission
    Array.prototype.filter.call(forms, function (form) {
      form.addEventListener('submit', function (event) {
        if (form.checkValidity() === false) {
          event.preventDefault();
          event.stopPropagation();
        }
        form.classList.add('was-validated');
      }, false);
    });
  }, false);
})();


function modalReservation(CSRF, roomID) {
  document.getElementById("check-availability-btn").addEventListener("click", function () {
    let modal_html = `
    <div class="container">
      <form id="availabilityModalForm" name="availabilityModalForm" action="" method="post" class="needs-validation" novalidate>
        <div class="row" id="reservation-dates-modal">
          <div class="col">
            <input required class="form-control" type="text" name="start_date" id="start_date" placeholder=" Arrival date" disabled>
          </div>
          <div class="col">
            <input required class="form-control" type="text" name="end_date" id="end_date" placeholder=" Departure date" disabled>
          </div>
        </div>
      </form>
    </div>`;
    attention.customModal({
      title: 'Choose your dates',
      msg: modal_html,
      inputAttributes: {},
      customClass: {},
      confirmButtonColor: "#0d6efd",
      willOpen: () => {
        const rdm = document.getElementById("reservation-dates-modal");
        const rp = new DateRangePicker(rdm, {
          format: 'dd/mm/yyyy',
          minDate: new Date(),
          todayHighlight: true,
          showOnFocus: true,
        })
      },
      didOpen: () => {
        document.getElementById("start_date").removeAttribute("disabled");
        document.getElementById("end_date").removeAttribute("disabled");
      },
      preConfirm: () => {
        return [
          document.getElementById('start_date').value,
          document.getElementById('end_date').value
        ]
      },
      callback: function(result) {
        let modalForm = document.getElementById("availabilityModalForm");
        let modalFormData = new FormData(modalForm);
        modalFormData.append("csrf_token", CSRF);
        modalFormData.append("room_id", roomID);

        fetch('/search-availability-modal', {
          method: "post",
          body: modalFormData,
        })
        .then(response => response.json())
        .then(data => {
          if (data.ok) {
            attention.customModal({
              title: "This room is available !",
              icon: "success",
              msg: '<p><a href="/reserve-room?id=' + data.room_id + '&sd=' + data.start_date + '&ed=' + data.end_date + '"' +
                ' class="btn btn-primary mt-4">Reserve Now !</a></p>',
              inputAttributes: {},
              customClass: {},
              showConfirmButton: false,
            })
          } else {
            attention.error({
              msg: "Sorry, this room is not available",
            });
          }
        });
      }
    });
  })
}

function notify(msg, msgType, duration) {
  notie.alert({
    type: msgType,
    text: msg,
    time: duration,
  })
}

function notifyModal(title, html, icon, confirmButtonText) {
  Swal.fire({
    title: title,
    html: html,
    icon: icon,
    confirmButtonText: confirmButtonText,
  })
}

function Inform() {
  let toast = function (params) {
    const {
      msg = '',
      icon = 'success',
      position = 'top-end',

    } = params;

    const Toast = Swal.mixin({
      toast: true,
      title: msg,
      position: position,
      icon: icon,
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
      didOpen: (toast) => {
        toast.addEventListener('mouseenter', Swal.stopTimer)
        toast.addEventListener('mouseleave', Swal.resumeTimer)
      }
    })

    Toast.fire({})
  }

  let success = function (params) {
    const {
      title = "",
      msg = "",
      footer = "",
    } = params;

    Swal.fire({
      icon: 'success',
      title: title,
      text: msg,
      footer: footer,
    })

  }

  let warning = function (params) {
    const {
      title = "",
      msg = "",
      footer = "",
    } = params;

    Swal.fire({
      icon: 'warning',
      title: title,
      text: msg,
      footer: footer,
    })

  }

  let error = function (params) {
    const {
      title = "",
      msg = "",
      footer = "",
    } = params;

    Swal.fire({
      icon: 'error',
      title: title,
      text: msg,
      footer: footer,
    })

  }

  let info = function (params) {
    const {
      title = "",
      msg = "",
      footer = "",
    } = params;

    Swal.fire({
      icon: 'info',
      title: title,
      text: msg,
      footer: footer,
    })

  }

  let question = function (params) {
    const {
      title = "",
      msg = "",
      footer = "",
    } = params;

    Swal.fire({
      icon: 'question',
      title: title,
      text: msg,
      footer: footer,
    })

  }

  async function customModal(params) {
    const {
      title = "",
      icon = "",
      msg = "",
      input = "",
      inputLabel = "",
      inputAttributes: {min: $min, max: $max, step: $step},
      inputValue: $inputValue,
      customClass: {sender: sender, content: htmlString},
      confirmButtonText = "OK",
      confirmButtonColor = "",
      showConfirmButton = true,
    } = params;

    const {value: result} = await Swal.fire({
      title: title,
      icon: icon,
      html: msg,
      input: input,
      inputLabel: inputLabel,
      inputAttributes: {min: $min, max: $max, step: $step},
      inputValue: $inputValue,
      customClass: {sender: sender, content: htmlString},
      backdrop: false,
      focusConfirm: false,
      confirmButtonText: confirmButtonText,
      confirmButtonColor: confirmButtonColor,
      showConfirmButton: showConfirmButton,
      showCancelButton: true,
      willOpen: () => {
        if (params.willOpen !== undefined) {
          params.willOpen();
        }
      },
      didOpen: () => {
        if (params.didOpen !== undefined) {
          params.didOpen();
        }
      },
      preConfirm: () => {
        if (params.preConfirm !== undefined) {
          params.preConfirm();
        }
      }
    })
    .then((result) => {
      if (result.dismiss === Swal.DismissReason.cancel) {
        // specific cancel action for blockDayRange() unchecking already checked status
        if (sender === 'blockDayRange()') {
          cbx = htmlString;
          cbx.checked = false;
        }
      }
      return result;
    });

    if (result) {
      if (result.dismiss !== Swal.DismissReason.cancel) {
        if (result.value !== "") {
          if (params.callback !== undefined) {
            params.callback(result);
          }
        } else {
          params.callback(false);
        }
      } else {
        params.callback(false);
      }
    }

  }

  return {
    toast: toast,
    success: success,
    warning: warning,
    error: error,
    info: info,
    question: question,
    customModal: customModal
  }
}
