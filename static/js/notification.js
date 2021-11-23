let attention = Prompt();

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

function notify(msg, msgType) {
  // set fixed navbar z-index = 0 allowing full noti(e)fication display
  document.getElementById("navbar").style.zIndex = 0
  notie.alert({
    type: msgType,
    text: msg,
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

function Prompt() {
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

  async function custom(params) {
    const {
      title = "",
      msg = "",
    } = params;

    const { value: formValues } = await Swal.fire({
      title: title,
      html: msg,
      backdrop: false,
      focusConfirm: false,
      showCancelButton: true,
      willOpen: () => {
        const rdm = document.getElementById("reservation-dates-modal");
        const rp = new DateRangePicker(rdm, {
          format: 'dd/mm/yyyy',
          showOnFocus: true,
        })
      },
      didOpen: () => {
        document.getElementById("start").removeAttribute("disabled");
        document.getElementById("end").removeAttribute("disabled");
      },
      preConfirm: () => {
        return [
          document.getElementById('start').value,
          document.getElementById('end').value
        ]
      }
    })

    if (formValues) {
      Swal.fire(JSON.stringify(formValues))
    }
  }

  return {
    toast: toast,
    success: success,
    error: error,
    custom: custom,
  }
}
