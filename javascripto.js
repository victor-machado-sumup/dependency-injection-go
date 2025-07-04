function calcularNetAmount(merchant, valorTransação) {
  let calculadoraTaxa;
  if (merchant.plano == "instant") {
    calculadoraTaxa = calcularTaxaInstant;
  } else if (merchant.plano == "economic") {
    calculadoraTaxa = calcularTaxaEconomic;
  } else if (merchant.plano == "accelerated") {
    calculadoraTaxa = calcularTaxaAccelerated;
  } else {
    throw new Error("plano não identificado");
  }

  return valorTransação - calculadoraTaxa(valorTransação);
}

function calcularTaxaInstant(valorTransação) {
  return valorTransação * 0.5;
}

function calcularTaxaEconomic(valorTransação) {
  return valorTransação * 0.1;
}

function calcularTaxaAccelerated(valorTransação) {
  return valorTransação * 0.2;
}

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//

function calcularNetAmount2(merchant, valorTransação) {
  const calculadora = obterCalculadora(merchant);
  return calcularNetAmountIoC(valorTransação, calculadora);
}

function calcularNetAmountIoC(valorTransação, calculadora) {
  return valorTransação - calculadora(valorTransação);
}

function obterCalculadora(merchant) {
  let calculadora;
  if (merchant.plano == "instant") {
    calculadora = calcularTaxaInstant;
  } else if (merchant.plano == "economic") {
    calculadora = calcularTaxaInstant;
  } else if (merchant.plano == "accelerated") {
    calculadora = calcularTaxaAccelerated;
  }
  return calculadora;
}

const merchant = {
  merchantCode: "M00000",
  plano: "instant",
};
const netAmount = calcularNetAmount2(merchant, 120);

console.log(netAmount);
