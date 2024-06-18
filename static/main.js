window.onload = function(){

document.getElementById("SubmitButton").onclick = function() {SendData()}

function SendData() {
   Subject = document.getElementById("SubjectField").value
    Email = document.getElementById("EmailField").value
    Body = document.getElementById("BodyField").value

    if( Subject.length === 0 || Email.length === 0 || Body.length === 0){
        alert('You need to fill all required fields!')
        return
    }
    httpinstance = new XMLHttpRequest()
    httpinstance.open( "POST", "http://localhost:5000/sendfourm", false )
    httpinstance.setRequestHeader("Content-Type", "application/json;charset=UTF-8")

    var data = JSON.stringify( {
        "Subject" : Subject,
        "Email" : Email,
        "Body" : Body
    })
    httpinstance.send(data)
    alert("Sent the data successfully")


}}
