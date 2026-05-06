export type FishQuality = "baik" | "sedang" | "buruk";

export type StockStatus = "available" | "depleted";

export type MovementType =
  | "in"
  | "out"
  | "quality_update"
  | "location_update"
  | "adjustment";

export type ApiResponse<T> = {
  success: boolean;
  message: string;
  data?: T;
  errors?: ApiFieldError[] | unknown;
};

export type ApiFieldError = {
  field: string;
  message: string;
};

export type FishType = {
  id: string;
  name: string;
  image_url?: string | null;
  description?: string | null;
  created_at: string;
  updated_at: string;
};

export type ColdStorage = {
  id: string;
  name: string;
  location_label?: string | null;
  description?: string | null;
  created_at: string;
  updated_at: string;
};

export type StockBatch = {
  id: string;
  fish_type_id: string;
  cold_storage_id: string;
  quality: FishQuality;
  initial_weight_kg: number;
  remaining_weight_kg: number;
  entered_at: string;
  status: StockStatus;
  notes?: string | null;
  created_at: string;
  updated_at: string;
};

export type FIFOStock = {
  id: string;
  fish_type_name: string;
  quality: FishQuality;
  remaining_weight_kg: number;
  entered_at: string;
  cold_storage_name: string;
  location_label?: string | null;
  fifo_rank: number;
};

export type DashboardSummary = {
  total_available_weight_kg: number;
  total_stock_batches: number;
  total_available_batches: number;
  total_depleted_batches: number;
  today_stock_in_kg: number;
  today_stock_out_kg: number;
  fish_type_summary: {
    fish_type_id: string;
    fish_type_name: string;
    available_weight_kg: number;
    available_batches: number;
  }[];
  cold_storage_summary: {
    cold_storage_id: string;
    cold_storage_name: string;
    available_weight_kg: number;
    available_batches: number;
  }[];
};

export type RecentMovement = {
  id: string;
  stock_batch_id: string;
  movement_type: MovementType;
  fish_type_name: string;
  weight_kg?: number | null;
  description: string;
  created_at: string;
};

export type StockOut = {
  id: string;
  fish_type_id?: string;
  destination: string;
  total_weight_kg: number;
  out_at: string;
  notes?: string | null;
  items: {
    stock_batch_id: string;
    fish_type_name?: string;
    weight_kg: number;
  }[];
  created_at: string;
};
