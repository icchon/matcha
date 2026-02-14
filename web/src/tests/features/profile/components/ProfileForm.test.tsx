import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ProfileForm } from '@/features/profile/components/ProfileForm';

const mockOnSubmit = vi.fn();

beforeEach(() => {
  vi.resetAllMocks();
});

describe('ProfileForm', () => {
  it('renders all required fields', () => {
    render(<ProfileForm onSubmit={mockOnSubmit} isLoading={false} />);

    expect(
      screen.getByLabelText(/first name/i),
      'ProfileForm should render a First Name input.',
    ).toBeInTheDocument();
    expect(
      screen.getByLabelText(/last name/i),
      'ProfileForm should render a Last Name input.',
    ).toBeInTheDocument();
    expect(
      screen.getByLabelText(/username/i),
      'ProfileForm should render a Username input.',
    ).toBeInTheDocument();
    expect(
      screen.getByLabelText(/gender/i),
      'ProfileForm should render a Gender select.',
    ).toBeInTheDocument();
    expect(
      screen.getByLabelText(/sexual preference/i),
      'ProfileForm should render a Sexual Preference select.',
    ).toBeInTheDocument();
    expect(
      screen.getByLabelText(/birthday/i),
      'ProfileForm should render a Birthday input.',
    ).toBeInTheDocument();
    expect(
      screen.getByLabelText(/biography/i),
      'ProfileForm should render a Biography textarea.',
    ).toBeInTheDocument();
  });

  it('renders submit button', () => {
    render(<ProfileForm onSubmit={mockOnSubmit} isLoading={false} />);

    expect(
      screen.getByRole('button', { name: /save/i }),
      'ProfileForm should have a Save button.',
    ).toBeInTheDocument();
  });

  it('shows validation errors for empty required fields', async () => {
    const user = userEvent.setup();
    render(<ProfileForm onSubmit={mockOnSubmit} isLoading={false} />);

    await user.click(screen.getByRole('button', { name: /save/i }));

    await waitFor(() => {
      expect(
        screen.getAllByRole('alert').length,
        'Validation errors should appear for empty required fields.',
      ).toBeGreaterThan(0);
    });
    expect(
      mockOnSubmit,
      'onSubmit should NOT be called when validation fails.',
    ).not.toHaveBeenCalled();
  });

  it('calls onSubmit with valid data', async () => {
    const user = userEvent.setup();
    render(<ProfileForm onSubmit={mockOnSubmit} isLoading={false} />);

    await user.type(screen.getByLabelText(/first name/i), 'John');
    await user.type(screen.getByLabelText(/last name/i), 'Doe');
    await user.type(screen.getByLabelText(/username/i), 'johndoe');
    await user.selectOptions(screen.getByLabelText(/gender/i), 'male');
    await user.selectOptions(screen.getByLabelText(/sexual preference/i), 'heterosexual');
    await user.type(screen.getByLabelText(/birthday/i), '1995-06-15');
    await user.type(screen.getByLabelText(/biography/i), 'Hello world');

    await user.click(screen.getByRole('button', { name: /save/i }));

    await waitFor(() => {
      expect(
        mockOnSubmit,
        'onSubmit should be called when all fields are valid.',
      ).toHaveBeenCalled();
    });

    const calledData = mockOnSubmit.mock.calls[0][0];
    expect(calledData.firstName, 'firstName should be John.').toBe('John');
    expect(calledData.lastName, 'lastName should be Doe.').toBe('Doe');
    expect(calledData.username, 'username should be johndoe.').toBe('johndoe');
    expect(calledData.gender, 'gender should be male.').toBe('male');
    expect(calledData.sexualPreference, 'sexualPreference should be heterosexual.').toBe('heterosexual');
    expect(calledData.birthday, 'birthday should be 1995-06-15.').toBe('1995-06-15');
    expect(calledData.biography, 'biography should be Hello world.').toBe('Hello world');
  });

  it('populates form with initial values', () => {
    const initialValues = {
      firstName: 'Jane',
      lastName: 'Smith',
      username: 'janesmith',
      gender: 'female' as const,
      sexualPreference: 'bisexual' as const,
      birthday: '1990-01-01',
      biography: 'Hi there',
    };
    render(<ProfileForm onSubmit={mockOnSubmit} isLoading={false} initialValues={initialValues} />);

    expect(
      (screen.getByLabelText(/first name/i) as HTMLInputElement).value,
      'First name should be pre-filled with initial value.',
    ).toBe('Jane');
    expect(
      (screen.getByLabelText(/biography/i) as HTMLTextAreaElement).value,
      'Biography should be pre-filled with initial value.',
    ).toBe('Hi there');
  });

  it('disables submit button when loading', () => {
    render(<ProfileForm onSubmit={mockOnSubmit} isLoading={true} />);

    expect(
      screen.getByRole('button', { name: /sav/i }),
      'Submit button should be disabled when isLoading is true.',
    ).toBeDisabled();
  });
});
