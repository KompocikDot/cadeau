export type Guest = {
  id: number;
  username: string;
};

export type UserOccasion = {
  id: number;
  name: string;
  guests: Guest[];
};

export type Gift = {
  id: number;
  name: string;
  url: string;
};
