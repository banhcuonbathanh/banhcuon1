import envConfig from "@/config";
import { AccountType } from "@/zusstand/account/domain/account.schema";


function delay(ms: number | undefined) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

// Example usage
async function fetchDataDelay() {
  console.log('Fetching data...');
  await delay(3000); // Delay for 30 seconds
  console.log('Data fetched!');
}



const get_Account = async (email : string): Promise<AccountType> => {
  try {
    const response = await fetch(
      envConfig.NEXT_PUBLIC_API_ENDPOINT +
      envConfig.NEXT_PUBLIC_API_Get_Account_Email + email,
      {
        method: "GET",
        cache: "no-store"
      }
    );

    const data = await response.json();
    console.log(
      "nextjs/app/(route)/product/controller-product/controller-product.ts",
      data.data
    );

    // Check if data is null
    if (data === null || data.data === null) {
      throw "data null ";
    }

    return data; // Returning the parsed JSON data
  } catch (error) {
    console.error("Error fetching billboard:", error);
    throw error; // Propagate the error to the caller
  }
};
export {
  fetchDataDelay,

  get_Account

};
