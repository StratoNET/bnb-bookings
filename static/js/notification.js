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

function notify(msg, msgType, duration) {
  // set fixed navbar z-index = 0 allowing full noti(e)fication display
  document.getElementById("navbar").style.zIndex = 0
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
      msg = "",
    } = params;

    const {value: result} = await Swal.fire({
      title: title,
      html: msg,
      backdrop: false,
      focusConfirm: false,
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
