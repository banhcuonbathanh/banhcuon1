import { SetInterface } from "@/schemaValidations/interface/types_set";
import SetCard from "./set";




interface SetCardListProps {
    sets: SetInterface[];
  }
  
  export function SetCardList({ sets }: SetCardListProps) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {sets.map((set) => (
          <SetCard key={set.id} set={set} />
        ))}
      </div>
    );
  }
  