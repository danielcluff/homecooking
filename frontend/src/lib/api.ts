const API_BASE_URL = 'http://localhost:8080/api/v1';

export interface Recipe {
	id: number;
	title: string;
	description?: string;
	slug: string;
	servings?: number;
	prep_time_minutes?: number;
	cook_time_minutes?: number;
	difficulty?: string;
	featured_image_path?: string;
	markdown_content: string;
	is_published: boolean;
	category_id?: number;
	owner_id: number;
	created_at: string;
	updated_at: string;
}

export interface Category {
	id: number;
	name: string;
	slug: string;
	created_at: string;
	updated_at: string;
}

export interface Tag {
	id: number;
	name: string;
	slug: string;
	color: string;
	created_at: string;
	updated_at: string;
}

class ApiClient {
	private baseUrl: string;

	constructor(baseUrl: string = API_BASE_URL) {
		this.baseUrl = baseUrl;
	}

	private getHeaders(includeAuth = true): HeadersInit {
		const headers: HeadersInit = {
			'Content-Type': 'application/json',
		};

		if (includeAuth) {
			const token = localStorage.getItem('accessToken');
			if (token) {
				headers['Authorization'] = `Bearer ${token}`;
			}
		}

		return headers;
	}

	async request<T>(
		endpoint: string,
		options: RequestInit = {}
	): Promise<T> {
		const url = `${this.baseUrl}${endpoint}`;
		const response = await fetch(url, {
			...options,
			headers: this.getHeaders(),
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({}));
			throw new Error(error.message || `Request failed: ${response.statusText}`);
		}

		return response.json();
	}

	async get<T>(endpoint: string): Promise<T> {
		return this.request<T>(endpoint, { method: 'GET' });
	}

	async post<T>(endpoint: string, data?: unknown): Promise<T> {
		return this.request<T>(endpoint, {
			method: 'POST',
			body: data ? JSON.stringify(data) : undefined,
		});
	}

	async put<T>(endpoint: string, data?: unknown): Promise<T> {
		return this.request<T>(endpoint, {
			method: 'PUT',
			body: data ? JSON.stringify(data) : undefined,
		});
	}

	async delete<T>(endpoint: string): Promise<T> {
		return this.request<T>(endpoint, { method: 'DELETE' });
	}
}

export const api = new ApiClient();

export async function getRecipe(id: number): Promise<Recipe> {
	return api.get<Recipe>(`/recipes/${id}`);
}

export async function getRecipes(params?: {
	limit?: number;
	offset?: number;
}): Promise<{ recipes: Recipe[]; total: number }> {
	const queryParams = new URLSearchParams();
	if (params?.limit) queryParams.append('limit', params.limit.toString());
	if (params?.offset) queryParams.append('offset', params.offset.toString());
	const query = queryParams.toString() ? `?${queryParams}` : '';
	return api.get<{ recipes: Recipe[]; total: number }>(`/recipes${query}`);
}

export async function createRecipe(data: Partial<Recipe>): Promise<Recipe> {
	return api.post<Recipe>('/recipes', data);
}

export async function updateRecipe(id: number, data: Partial<Recipe>): Promise<Recipe> {
	return api.put<Recipe>(`/recipes/${id}`, data);
}

export async function deleteRecipe(id: number): Promise<void> {
	return api.delete<void>(`/recipes/${id}`);
}

export async function getCategories(): Promise<Category[]> {
	return api.get<Category[]>('/categories');
}

export async function getCategory(id: number): Promise<Category> {
	return api.get<Category>(`/categories/${id}`);
}

export async function createCategory(data: Partial<Category>): Promise<Category> {
	return api.post<Category>('/categories', data);
}

export async function updateCategory(id: number, data: Partial<Category>): Promise<Category> {
	return api.put<Category>(`/categories/${id}`, data);
}

export async function deleteCategory(id: number): Promise<void> {
	return api.delete<void>(`/categories/${id}`);
}

export async function getTags(): Promise<Tag[]> {
	return api.get<Tag[]>('/tags');
}

export async function getTag(id: number): Promise<Tag> {
	return api.get<Tag>(`/tags/${id}`);
}

export async function createTag(data: Partial<Tag>): Promise<Tag> {
	return api.post<Tag>('/tags', data);
}

export async function updateTag(id: number, data: Partial<Tag>): Promise<Tag> {
	return api.put<Tag>(`/tags/${id}`, data);
}

export async function deleteTag(id: number): Promise<void> {
	return api.delete<void>(`/tags/${id}`);
}

export async function apiFetch<T>(endpoint: string, authenticated = false): Promise<T> {
	const headers: HeadersInit = {
		'Content-Type': 'application/json',
	};

	if (authenticated) {
		const token = localStorage.getItem('accessToken');
		if (token) {
			headers['Authorization'] = `Bearer ${token}`;
		}
	}

	const response = await fetch(`${API_BASE_URL}${endpoint}`, {
		headers,
	});

	if (!response.ok) {
		const error = await response.json().catch(() => ({}));
		throw new Error(error.message || `Request failed: ${response.statusText}`);
	}

	return response.json();
}
