import ReadingTestPage from "./(landingpage)/page";
import Home from "./(landingpage)/page";

export default async function DashboardLayout({
  children
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="mx-auto max-w-7xl ">
      <p className="">lay out</p>

      <ReadingTestPage />
    </div>
  );
}
