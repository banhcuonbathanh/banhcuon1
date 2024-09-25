

function delay(ms: number | undefined) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

// Example usage
async function fetchDataDelay() {
  console.log('Fetching data...');
  await delay(3000); // Delay for 30 seconds
  console.log('Data fetched!');
}

export {
  fetchDataDelay

};
