if(document.getElementById("syes").checked){
    showSmokeCount(0)
}
else{
    hideSmokeCount(0)
}

let bmi=0
calcBMI()

function hideSmokeCount(that) {
    var ele=document.getElementById('cpd')
    lbl=document.getElementById('labelcpd')
    ele.style.display="none"
    lbl.style.display="none"

    lbl.float

    ele.value=0
  }
function showSmokeCount(that){
    lbl=document.getElementById('labelcpd')
    var ele=document.getElementById('cpd')
    lbl.style.display="block"
    ele.style.display="block"
    ele.value=1
  }
function calcBMI(){
    h=document.getElementById('height').value
    w=document.getElementById('weight').value

    let h_m=h/100.0
    bmi=w/(h_m*h_m)
    document.getElementById("bmi").innerHTML="BMI= "+bmi.toFixed(2)
}