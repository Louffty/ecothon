import React, { useEffect, useState } from 'react';

const SwitchMoney = () => {
  const [isToggled, setIsToggled] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('http://nothypeproduction.space/api/v1/admin/all');
        if (!response.ok) throw new Error('Network response was not ok');
        const data = await response.json();
        setIsToggled(data.body.value);
      } catch (error) {
        console.error('Error fetching data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const handleToggle = async () => {
    const newValue = !isToggled;
    setIsToggled(newValue);
    try {
      const response = await fetch('https://nothypeproduction.space/api/v1/admin/update', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          value: newValue,
          uuid: '07da62f6-fc0a-4ae3-964a-2176442afee8'
        })
      });

      if (!response.ok) throw new Error('Network response was not ok');
    } catch (error) {
      console.error('Error updating toggle:', error);
      console.log(newValue)
      setIsToggled(!newValue);
    }
  };

  if (loading) {
    return <p>Loading...</p>;
  }

  return (
    <div>
      <input
        type="checkbox"
        id="toggle"
        checked={isToggled}
        onChange={handleToggle}
      />
      <label htmlFor="toggle">
        <span></span>
      </label>
    </div>
  );
};

export default SwitchMoney;