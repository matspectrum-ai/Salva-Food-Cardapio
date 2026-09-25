import { streamToAPI } from "@/lib/bff";

export async function GET() {
  return streamToAPI("/api/v1/orders/events");
}
