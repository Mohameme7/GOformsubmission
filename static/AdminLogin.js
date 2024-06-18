window.onload =  function(){
  loginbutton = document.getElementById("loginbutton").onclick = function () {submitLogin()}

  function submitLogin() {

    passwordvalue = document.getElementById("passwordinput").value
    var data = {
        Password: passwordvalue
    };
    if(passwordvalue.length === 0){
      alert("Please write a password")
      return
    }

    fetch('/admin/login', {
        method: 'POST',
        body: JSON.stringify(data)
    })
     .then(response => {
        if (!response.ok) {
            alert('Wrong Password');
            throw new Error('Login failed');

        }
    })
    .then(data => {
        console.log('Login successful:', data);
        window.location.replace("/panel");
    })}}




